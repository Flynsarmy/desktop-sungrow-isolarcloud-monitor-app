package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	stdruntime "runtime"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	credentials   *Credentials
	httpClient    *http.Client
	TrayTitleChan chan string
	TrayIconChan  chan []byte
	BaseIcon      []byte
}

// Credentials stores API authentication data
type Credentials struct {
	AppKey       string `json:"appKey"`
	SecretKey    string `json:"secretKey"`
	AuthURL      string `json:"authUrl"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	TokenExpiry  int64  `json:"tokenExpiry,omitempty"`
	GatewayURL   string `json:"gatewayUrl,omitempty"`
}

// Plant represents a solar plant
type Plant struct {
	PsID                 int     `json:"ps_id"`
	PsName               string  `json:"ps_name"`
	Description          *string `json:"description"`
	PsType               int     `json:"ps_type"`
	OnlineStatus         int     `json:"online_status"`
	ValidFlag            int     `json:"valid_flag"`
	GridConnectionStatus int     `json:"grid_connection_status"`
	InstallDate          string  `json:"install_date"`
	PsLocation           string  `json:"ps_location"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	PsFaultStatus        int     `json:"ps_fault_status"`
	ConnectType          int     `json:"connect_type"`
	UpdateTime           string  `json:"update_time"`
	PsCurrentTimeZone    string  `json:"ps_current_time_zone"`
	GridConnectionTime   *string `json:"grid_connection_time"`
	BuildStatus          int     `json:"build_status"`
	TodayEnergy          string  `json:"today_energy,omitempty"`
}

// PlantDevice represents a device in a plant
type PlantDevice struct {
	UUID               int    `json:"uuid"`
	PsKey              string `json:"ps_key"`
	DeviceSN           string `json:"device_sn"`
	DeviceName         string `json:"device_name"`
	DeviceType         int    `json:"device_type"`
	TypeName           string `json:"type_name"`
	DeviceModelID      int    `json:"device_model_id"`
	DeviceModelCode    string `json:"device_model_code"`
	DevFaultStatus     int    `json:"dev_fault_status"`
	DevStatus          string `json:"dev_status"`
	ClaimState         int    `json:"claim_state"`
	DeviceCode         int    `json:"device_code"`
	ChnnlID            int    `json:"chnnl_id"`
	CommunicationDevSN string `json:"communication_dev_sn"`
	PsID               int    `json:"ps_id"`
}

// ApiResponse wraps API responses
type ApiResponse struct {
	ReqSerialNum string          `json:"req_serial_num,omitempty"`
	ResultCode   string          `json:"result_code"`
	ResultMsg    string          `json:"result_msg"`
	ResultData   json.RawMessage `json:"result_data"`
}

// LoginResultData contains OAuth token data
type LoginResultData struct {
	AccessToken  string   `json:"access_token"`
	TokenType    string   `json:"token_type"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	AuthPsList   []string `json:"auth_ps_list"`
	AuthUser     int      `json:"auth_user"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		TrayTitleChan: make(chan string, 10),
		TrayIconChan:  make(chan []byte, 10),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadCredentials()
}

// GetStoredCredentials returns stored credentials
func (a *App) GetStoredCredentials() *Credentials {
	if a.credentials == nil {
		a.loadCredentials()
	}
	return a.credentials
}

// Authenticate handles OAuth flow
func (a *App) Authenticate(creds Credentials) (map[string]interface{}, error) {
	// Store credentials
	a.credentials = &creds

	// Find an available port
	port := 0
	var server *http.Server
	var listener net.Listener

	for p := 8080; p <= 8090; p++ {
		addr := fmt.Sprintf(":%d", p)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			port = p
			listener = l
			server = &http.Server{}
			break
		}
	}

	if port == 0 {
		return nil, fmt.Errorf("no available ports found (tried 8080-8090)")
	}

	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)

	// Build auth URL
	authURL, err := url.Parse(creds.AuthURL)
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("invalid auth URL: %w", err)
	}

	query := authURL.Query()
	query.Set("redirectUrl", redirectURL)
	authURL.RawQuery = query.Encode()

	// Create a channel to receive the authorization code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Setup HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Authentication failed: no code received"))
			return
		}

		// Send success page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
			<head><title>Authentication Successful</title></head>
			<body style="font-family: sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); color: white;">
				<div style="text-align: center;">
					<h1 style="color: #10b981;">✓ Authentication Successful</h1>
					<p>You can close this window and return to the app.</p>
				</div>
			</body>
			</html>
		`))

		codeChan <- code
	})

	server.Handler = mux

	// Start server in background
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Open browser for OAuth
	runtime.BrowserOpenURL(a.ctx, authURL.String())

	// Wait for code or error (with timeout)
	select {
	case code := <-codeChan:
		// Shutdown server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)

		// Exchange code for tokens
		return a.exchangeCodeForTokens(code, creds, redirectURL)

	case err := <-errChan:
		server.Shutdown(context.Background())
		return nil, err

	case <-time.After(5 * time.Minute):
		server.Shutdown(context.Background())
		return nil, fmt.Errorf("authentication timeout")
	}
}

// exchangeCodeForTokens exchanges authorization code for access tokens
func (a *App) exchangeCodeForTokens(code string, creds Credentials, redirectURL string) (map[string]interface{}, error) {
	gatewayURL := creds.GatewayURL
	if gatewayURL == "" {
		gatewayURL = "https://augateway.isolarcloud.com"
	}

	tokenURL := fmt.Sprintf("%s/openapi/apiManage/token", gatewayURL)

	reqBody := map[string]interface{}{
		"appkey":       creds.AppKey,
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": redirectURL,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-access-key", creds.SecretKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.ResultCode != "1" {
		return nil, fmt.Errorf("authentication failed: %s", apiResp.ResultMsg)
	}

	var loginData LoginResultData
	if err := json.Unmarshal(apiResp.ResultData, &loginData); err != nil {
		return nil, err
	}

	// Calculate expiry time
	expiry := time.Now().Add(time.Duration(loginData.ExpiresIn) * time.Second).Unix()

	// Update credentials with tokens
	a.credentials.AccessToken = loginData.AccessToken
	a.credentials.RefreshToken = loginData.RefreshToken
	a.credentials.TokenExpiry = expiry * 1000 // Convert to milliseconds

	// Save credentials
	if err := a.saveCredentials(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"authenticated": true,
		"tokenExpiry":   a.credentials.TokenExpiry,
	}, nil
}

// GetPlantList retrieves list of solar plants
func (a *App) GetPlantList() ([]Plant, error) {
	if a.credentials == nil || a.credentials.AccessToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	gatewayURL := a.credentials.GatewayURL
	if gatewayURL == "" {
		gatewayURL = "https://augateway.isolarcloud.com"
	}

	apiURL := fmt.Sprintf("%s/openapi/platform/queryPowerStationList", gatewayURL)
	fmt.Printf("GetPlantList: Calling %s\n", apiURL)

	reqBody := map[string]interface{}{
		"appkey": a.credentials.AppKey,
		"page":   1,
		"size":   50,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	fmt.Printf("GetPlantList: Request body: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.credentials.AccessToken)
	req.Header.Set("x-access-key", a.credentials.SecretKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		fmt.Printf("GetPlantList: HTTP error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("GetPlantList: Response status: %d, body: %s\n", resp.StatusCode, string(body))

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("GetPlantList: JSON unmarshal error: %v\n", err)
		return nil, err
	}

	if apiResp.ResultCode != "1" {
		fmt.Printf("GetPlantList: API error code: %s, message: %s\n", apiResp.ResultCode, apiResp.ResultMsg)
		return nil, fmt.Errorf("API error: %s", apiResp.ResultMsg)
	}

	var result struct {
		PageList []Plant `json:"pageList"`
	}
	if err := json.Unmarshal(apiResp.ResultData, &result); err != nil {
		fmt.Printf("GetPlantList: Result data unmarshal error: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetPlantList: Successfully loaded %d plants\n", len(result.PageList))
	return result.PageList, nil
}

// GetDeviceList retrieves devices for a plant
func (a *App) GetDeviceList(psID int) ([]PlantDevice, error) {
	if a.credentials == nil || a.credentials.AccessToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	gatewayURL := a.credentials.GatewayURL
	if gatewayURL == "" {
		gatewayURL = "https://augateway.isolarcloud.com"
	}

	apiURL := fmt.Sprintf("%s/openapi/platform/getDeviceListByPsId", gatewayURL)

	reqBody := map[string]interface{}{
		"appkey": a.credentials.AppKey,
		"ps_id":  fmt.Sprintf("%d", psID),
		"page":   1,
		"size":   50,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.credentials.AccessToken)
	req.Header.Set("x-access-key", a.credentials.SecretKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.ResultCode != "1" {
		return nil, fmt.Errorf("API error: %s", apiResp.ResultMsg)
	}

	var result struct {
		PageList []PlantDevice `json:"pageList"`
	}
	if err := json.Unmarshal(apiResp.ResultData, &result); err != nil {
		return nil, err
	}

	return result.PageList, nil
}

// GetDevicePointData retrieves real-time data points for a device
func (a *App) GetDevicePointData(deviceType int, psKey string, pointIDs []int) ([]map[string]interface{}, error) {
	if a.credentials == nil || a.credentials.AccessToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	gatewayURL := a.credentials.GatewayURL
	if gatewayURL == "" {
		gatewayURL = "https://augateway.isolarcloud.com"
	}

	apiURL := fmt.Sprintf("%s/openapi/platform/getDeviceRealTimeData", gatewayURL)

	// Convert point IDs to strings
	pointIDStrs := make([]string, len(pointIDs))
	for i, id := range pointIDs {
		pointIDStrs[i] = fmt.Sprintf("%d", id)
	}

	reqBody := map[string]interface{}{
		"appkey":            a.credentials.AppKey,
		"device_type":       deviceType,
		"ps_key_list":       []string{psKey},
		"point_id_list":     pointIDStrs,
		"is_get_point_dict": "1",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.credentials.AccessToken)
	req.Header.Set("x-access-key", a.credentials.SecretKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.ResultCode != "1" {
		return nil, fmt.Errorf("API error: %s", apiResp.ResultMsg)
	}

	var result struct {
		DevicePointList []struct {
			DevicePoint map[string]interface{} `json:"device_point"`
		} `json:"device_point_list"`
	}
	if err := json.Unmarshal(apiResp.ResultData, &result); err != nil {
		return nil, err
	}

	// Extract device points
	devicePoints := make([]map[string]interface{}, len(result.DevicePointList))
	for i, item := range result.DevicePointList {
		devicePoints[i] = item.DevicePoint
	}

	return devicePoints, nil
}

// Logout clears stored credentials
func (a *App) Logout() error {
	a.credentials = nil
	return a.saveCredentials()
}

// loadCredentials loads credentials from file
func (a *App) loadCredentials() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	appDir := filepath.Join(configDir, "SungrowMonitor")
	credFile := filepath.Join(appDir, "credentials.json")

	data, err := os.ReadFile(credFile)
	if err != nil {
		return
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return
	}

	a.credentials = &creds
}

// saveCredentials saves credentials to file
func (a *App) saveCredentials() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, "SungrowMonitor")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	credFile := filepath.Join(appDir, "credentials.json")

	if a.credentials == nil {
		// Delete file if credentials are nil
		os.Remove(credFile)
		return nil
	}

	data, err := json.MarshalIndent(a.credentials, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(credFile, data, 0600)
}

// UpdateTrayTitle updates the system tray title/tooltip
func (a *App) UpdateTrayTitle(title string) {
	select {
	case a.TrayTitleChan <- title:
	default:
		// Channel full, skip update
	}
}

// UpdateTrayStatus updates the system tray icon with battery percentage and title
func (a *App) UpdateTrayStatus(percentage int, title string) {
	fmt.Printf("UpdateTrayStatus called: %d%%, %s\n", percentage, title)

	// Update Title
	a.UpdateTrayTitle(title)

	// Update Icon
	iconBytes, err := a.generateIconWithBadge(percentage)
	if err != nil {
		fmt.Printf("Failed to generate icon: %v\n", err)
		return
	}
	fmt.Printf("Generated icon bytes: %d\n", len(iconBytes))

	select {
	case a.TrayIconChan <- iconBytes:
		fmt.Println("Sent icon to TrayIconChan")
	default:
		fmt.Println("TrayIconChan blocked/full")
	}
}

func (a *App) generateIconWithBadge(percentage int) ([]byte, error) {
	width := 32
	height := 32
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// Colors
	bgColor := color.RGBA{80, 80, 80, 255} // Dark Gray background for empty part
	var fgColor color.RGBA

	if percentage <= 20 {
		fgColor = color.RGBA{220, 38, 38, 255} // Red
	} else if percentage <= 50 {
		fgColor = color.RGBA{234, 179, 8, 255} // Yellow/Orange
	} else {
		fgColor = color.RGBA{22, 163, 74, 255} // Green
	}

	// Geometry
	cx := float64(width) / 2
	cy := float64(height) / 2
	r := 15.0 // Radius

	// Angle limit (percentage 0..100 -> 0..2*Pi)
	limit := (float64(percentage) / 100.0) * 2 * math.Pi

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x) - cx + 0.5 // +0.5 for pixel center
			dy := float64(y) - cy + 0.5
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= r {
				// Inside circle
				// Calculate angle from Top (-Pi/2) Clockwise
				angle := math.Atan2(dy, dx) // -Pi to Pi

				// Standard Atan2: Right=0, Down=PI/2, Left=PI, Top=-PI/2
				// We want Top=0, Right=PI/2, Down=PI, Left=3PI/2

				angle += math.Pi / 2
				if angle < 0 {
					angle += 2 * math.Pi
				}

				if angle <= limit {
					rgba.Set(x, y, fgColor)
				} else {
					rgba.Set(x, y, bgColor)
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	err := png.Encode(buf, rgba)
	if err != nil {
		return nil, err
	}

	// Windows requires usage of ICO format for system tray
	if stdruntime.GOOS == "windows" {
		return a.convertToIco(buf.Bytes())
	}

	// macOS and Linux generally prefer PNG
	return buf.Bytes(), nil
}

func (a *App) convertToIco(pngData []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// ICONDIR header
	binary.Write(buf, binary.LittleEndian, uint16(0)) // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Type 1 = Icon
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Count = 1

	// ICONDIRENTRY
	buf.WriteByte(32)                                            // Width
	buf.WriteByte(32)                                            // Height
	buf.WriteByte(0)                                             // ColorCount (0 for >= 8bpp)
	buf.WriteByte(0)                                             // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1))            // Planes
	binary.Write(buf, binary.LittleEndian, uint16(32))           // BitCount
	binary.Write(buf, binary.LittleEndian, uint32(len(pngData))) // SizeInBytes
	binary.Write(buf, binary.LittleEndian, uint32(22))           // Offset (6 + 16)

	// PNG Data
	buf.Write(pngData)

	return buf.Bytes(), nil
}
