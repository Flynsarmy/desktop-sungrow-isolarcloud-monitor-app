import React, { useState, useEffect } from 'react'
import { PlantDevice } from './PlantDevice'
import { PlantDeviceBattery } from './PlantDeviceBattery'
import { Layers } from 'lucide-react'
import { GetDeviceList } from '../../wailsjs/go/main/App'

interface PlantDeviceType {
    uuid: number
    ps_key: string
    device_sn: string
    device_name: string
    device_type: number
    type_name: string
    device_model_id: number
    device_model_code: string
    dev_fault_status: number
    dev_status: string
    claim_state: number
    device_code: number
    chnnl_id: number
    communication_dev_sn: string
    ps_id: number
}

interface PlantDeviceListProps {
    ps_id: number
}

export function PlantDeviceList({ ps_id }: PlantDeviceListProps) {
    const [devices, setDevices] = useState<PlantDeviceType[]>([])
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (ps_id) {
            loadDevices()
        }
    }, [ps_id])

    const loadDevices = async () => {
        setIsLoading(true)
        setError(null)
        try {
            const deviceList = await GetDeviceList(ps_id)
            setDevices(deviceList || [])
        } catch (e: any) {
            console.error('Failed to load devices', e)
            setError(e.message || 'Failed to load devices')
        } finally {
            setIsLoading(false)
        }
    }

    if (isLoading) {
        return <div className="loading-compact">Finding devices...</div>
    }

    if (error) {
        return <div className="error-compact">{error}</div>
    }

    return (
        <div className="device-list-container">
            <h3 className="section-title">
                <Layers size={14} />
                Devices ({devices.length})
            </h3>

            {devices.length === 0 ? (
                <div className="empty-state">
                    <p>No devices found for this plant.</p>
                </div>
            ) : (
                <div className="device-grid">
                    {devices.map((device) =>
                        device.device_type === 43 ? (
                            <PlantDeviceBattery key={device.uuid || device.device_sn} device={device} />
                        ) : (
                            <PlantDevice key={device.uuid || device.device_sn} device={device} />
                        )
                    )}
                </div>
            )}
        </div>
    )
}
