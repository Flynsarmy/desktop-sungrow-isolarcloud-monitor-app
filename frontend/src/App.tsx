import React, { useState, useEffect } from 'react'
import { Login } from './components/Login'
import { PlantDetails } from './components/PlantDetails'
import { ArrowLeft } from 'lucide-react'
import { GetStoredCredentials, GetPlantList, Authenticate, Logout } from '../wailsjs/go/main/App'

function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false)
    const [isLoading, setIsLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [plants, setPlants] = useState<any[]>([])
    const [selectedPlant, setSelectedPlant] = useState<any | null>(null)

    useEffect(() => {
        checkAuth()
    }, [])

    const checkAuth = async () => {
        setIsLoading(true)
        try {
            const creds = await GetStoredCredentials()
            if (creds && creds.accessToken && creds.tokenExpiry && creds.tokenExpiry > Date.now()) {
                setIsAuthenticated(true)
                loadPlants()
            }
        } catch (err) {
            console.error('Auth check failed:', err)
        } finally {
            setIsLoading(false)
        }
    }

    const loadPlants = async () => {
        try {
            const plantList = await GetPlantList()
            console.log('Plants loaded:', plantList)
            setPlants(plantList || [])
        } catch (err: any) {
            console.error('Failed to load plants:', err)
            const errorMsg = err?.message || err?.toString() || 'Unknown error occurred'
            setError('Failed to load plants: ' + errorMsg)
        }
    }

    const handleLogin = async (credentials: any) => {
        setIsLoading(true)
        setError(null)
        try {
            const result = await Authenticate(credentials)
            if (result.authenticated) {
                setIsAuthenticated(true)
                await loadPlants()
            } else {
                setError(result.message || 'Authentication pending')
            }
        } catch (err: any) {
            setError(err.message || 'Authentication failed')
        } finally {
            setIsLoading(false)
        }
    }

    const handleLogout = async () => {
        await Logout()
        setIsAuthenticated(false)
        setPlants([])
        setSelectedPlant(null)
    }

    if (isLoading && !isAuthenticated) {
        return (
            <div className="container" style={{ justifyContent: 'center', alignItems: 'center' }}>
                <div className="status-badge">Loading...</div>
            </div>
        )
    }

    return (
        <div className="container">
            <header>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    {selectedPlant && (
                        <button
                            onClick={() => setSelectedPlant(null)}
                            style={{
                                background: 'none',
                                padding: '0.5rem',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center'
                            }}
                        >
                            <ArrowLeft size={20} />
                        </button>
                    )}
                    <h1>{selectedPlant ? selectedPlant.ps_name : 'Sungrow iSolarCloud'}</h1>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    {isAuthenticated && (
                        <button onClick={handleLogout} style={{ padding: '0.25rem 0.75rem', fontSize: '0.75rem' }}>
                            Logout
                        </button>
                    )}
                    <div
                        className={`status-badge ${isAuthenticated ? '' : 'error'}`}
                        style={
                            !isAuthenticated
                                ? {
                                    color: '#ef4444',
                                    backgroundColor: 'rgba(239, 68, 68, 0.1)',
                                    borderColor: 'rgba(239, 68, 68, 0.2)'
                                }
                                : {}
                        }
                    >
                        {isAuthenticated ? 'Connected' : 'Disconnected'}
                    </div>
                </div>
            </header>

            <main>
                {error && (
                    <div
                        className="card"
                        style={{ borderLeft: '4px solid #ef4444', marginBottom: '1.5rem', padding: '1rem' }}
                    >
                        <p style={{ margin: 0, color: '#ef4444', fontSize: '0.875rem' }}>{error}</p>
                    </div>
                )}

                {!isAuthenticated ? (
                    <Login onLogin={handleLogin} isLoading={isLoading} />
                ) : selectedPlant ? (
                    <PlantDetails plant={selectedPlant} />
                ) : (
                    <div
                        style={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
                            gap: '1.5rem'
                        }}
                    >
                        {plants.map((plant) => (
                            <div key={plant.ps_id} className="card">
                                <h3 style={{ margin: '0 0 1rem 0' }}>{plant.ps_name}</h3>
                                <div style={{ fontSize: '0.875rem', color: '#94a3b8' }}>
                                    <p>Location: {plant.ps_location}</p>
                                    <p>Status: {plant.ps_fault_status === 3 ? 'Normal' : 'Attention'}</p>
                                    <p>Daily Yield: {plant.today_energy || '0'} kWh</p>
                                </div>
                                <button
                                    style={{ width: '100%', marginTop: '1rem' }}
                                    onClick={() => setSelectedPlant(plant)}
                                >
                                    View Details
                                </button>
                            </div>
                        ))}
                        {plants.length === 0 && <p>No plants found.</p>}
                    </div>
                )}
            </main>
        </div>
    )
}

export default App
