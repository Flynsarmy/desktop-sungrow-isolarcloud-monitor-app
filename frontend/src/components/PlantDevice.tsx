import React from 'react'
import {
    Zap,
    Power,
    Boxes,
    CloudSun,
    Gauge,
    Cpu,
    Factory,
    Home,
    Container,
    Settings,
    Battery,
    Fuel,
    Sun,
    Box
} from 'lucide-react'

interface PlantDeviceType {
    device_type: number
    device_name: string
    dev_fault_status: number
    type_name: string
    ps_key: string
}

const DEVICE_ICONS: Record<number, any> = {
    1: Zap, // Inverter
    3: Power, // Grid-Connection Point
    4: Boxes, // Combiner Box
    5: CloudSun, // Meteo Station
    7: Gauge, // Meter
    9: Cpu, // Data Logger
    11: Factory, // Plant
    14: Home, // Energy Storage System
    17: Container, // Unit
    41: Settings, // Optimizer
    43: Battery, // Battery
    51: Fuel, // Charger
    55: Sun // Microinverter
}

interface PlantDeviceProps {
    device: PlantDeviceType
}

export function PlantDevice({ device }: PlantDeviceProps) {
    const Icon = DEVICE_ICONS[device.device_type] || Box

    return (
        <div className="device-card">
            <div className="device-icon">
                <Icon size={18} />
            </div>
            <div className="device-info">
                <div className="device-header">
                    <span className="device-name">{device.device_name}</span>
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
