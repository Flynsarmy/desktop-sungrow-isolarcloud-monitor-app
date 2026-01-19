import React, { useState, useEffect } from 'react'
import { Battery } from 'lucide-react'
import { GetDevicePointData, UpdateTrayStatus } from '../../wailsjs/go/main/App'

interface PlantDeviceType {
    device_type: number
    device_name: string
    dev_fault_status: number
    type_name: string
    ps_key: string
}

interface PlantDeviceBatteryProps {
    device: PlantDeviceType
}

export function PlantDeviceBattery({ device }: PlantDeviceBatteryProps) {
    const [soc, setSoc] = useState<number | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function fetchSoc() {
            try {
                const data = await GetDevicePointData(device.device_type, device.ps_key, [58604])

                if (data && data.length > 0) {
                    const pointData = data[0]
                    const socValue = pointData['p58604']
                    if (socValue !== undefined && socValue !== null) {
                        const val = Math.round(parseFloat(socValue) * 1000) / 10
                        setSoc(val)
                        // Update tray status with battery percentage
                        const trayTitle = `Battery: ${Math.round(val)}%`
                        console.log('Updating tray status:', trayTitle)
                        UpdateTrayStatus(Math.round(val), trayTitle)
                    }
                }
            } catch (error) {
                console.error('Failed to fetch battery SOC:', error)
            } finally {
                setLoading(false)
            }
        }

        fetchSoc()

        // Refresh every 5 minutes
        const interval = setInterval(fetchSoc, 5 * 60 * 1000)

        // Cleanup interval on unmount
        return () => clearInterval(interval)
    }, [device.ps_key, device.device_type])

    return (
        <div className="device-card battery">
            <div className="device-icon">
                <Battery size={18} />
            </div>
            <div className="device-info">
                <div className="device-header">
                    <span className="device-name">{device.device_name}</span>
                    {loading ? (
                        <span className="loading-text">Loading...</span>
                    ) : soc !== null ? (
                        <span className="soc-value">{soc}%</span>
                    ) : null}
                    <span className={`device-status ${device.dev_fault_status === 4 ? 'normal' : 'fault'}`}>
                        {device.dev_fault_status === 4 ? 'Normal' : 'Fault'}
                    </span>
                </div>
                <div className="device-meta">
                    <span>{device.type_name}</span>
                    <span className="separator">|</span>
                    <span className="mono">{device.ps_key}</span>
                </div>
            </div>
        </div>
    )
}
