import React, { useState } from 'react'

interface LoginProps {
    onLogin: (creds: any) => void
    isLoading: boolean
}

export function Login({ onLogin, isLoading }: LoginProps) {
    const [appKey, setAppKey] = useState('')
    const [secretKey, setSecretKey] = useState('')
    const [authUrl, setAuthUrl] = useState('https://auapi.isolarcloud.com:443/openapi/apiManage/token')
    const [gatewayUrl, setGatewayUrl] = useState('https://augateway.isolarcloud.com')

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onLogin({ appKey, secretKey, authUrl, gatewayUrl })
    }

    return (
        <div className="card" style={{ maxWidth: '400px', margin: '2rem auto' }}>
            <h2 style={{ marginBottom: '1.5rem', textAlign: 'center' }}>Connect to Sungrow</h2>
            <form onSubmit={handleSubmit}>
                <div className="input-group">
                    <label>App Key</label>
                    <input
                        type="text"
                        value={appKey}
                        onChange={(e) => setAppKey(e.target.value)}
                        required
                        placeholder="Enter your App Key"
                    />
                </div>
                <div className="input-group">
                    <label>Secret Key</label>
                    <input
                        type="password"
                        value={secretKey}
                        onChange={(e) => setSecretKey(e.target.value)}
                        required
                        placeholder="Enter your Secret Key"
                    />
                </div>
                <div className="input-group">
                    <label>Authorization URL</label>
                    <input
                        type="text"
                        value={authUrl}
                        onChange={(e) => setAuthUrl(e.target.value)}
                        required
                    />
                </div>
                <div className="input-group">
                    <label>Country / Gateway</label>
                    <select
                        value={gatewayUrl}
                        onChange={(e) => setGatewayUrl(e.target.value)}
                        required
                    >
                        <option value="https://augateway.isolarcloud.com">Australia</option>
                        <option value="https://gateway.isolarcloud.com">China</option>
                        <option value="https://gateway.isolarcloud.com.hk">International</option>
                        <option value="https://gateway.isolarcloud.eu">Europe</option>
                    </select>
                </div>
                <button type="submit" style={{ width: '100%', marginTop: '1rem' }} disabled={isLoading}>
                    {isLoading ? 'Authenticating...' : 'Authenticate'}
                </button>
            </form>
            <p style={{ marginTop: '1.5rem', fontSize: '0.75rem', color: '#94a3b8', textAlign: 'center' }}>
                Get your API credentials from the <a href="https://developer-api.isolarcloud.com" target="_blank" rel="noreferrer" style={{ color: '#f59e0b' }}>Sungrow Developer Portal</a>
            </p>
        </div>
    )
}
