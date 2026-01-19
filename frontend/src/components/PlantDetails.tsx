import React from 'react'
import { PlantDeviceList } from './PlantDeviceList'

const PLANT_TYPES: Record<number, string> = {
    1: 'Utility Plant',
    3: 'Distributed PV',
    4: 'Residential PV',
    5: 'Residential Storage',
    6: 'Village Plant',
    7: 'Dist. Storage',
    8: 'Poverty Alleviation',
    9: 'Wind Power',
    12: 'C&I Storage'
}

interface Plant {
    ps_id: number
    ps_name: string
    ps_location: string
    ps_type: number
    ps_fault_status: number
    online_status: number
    install_date?: string
}

interface PlantDetailsProps {
    plant: Plant
}

export function PlantDetails({ plant }: PlantDetailsProps) {
    return (
        <div className="plant-details-view">
            <div className="plant-info-grid card">
                <DetailRow label="ID" value={String(plant.ps_id)} />
                <DetailRow label="Name" value={plant.ps_name} />
                <DetailRow label="Location" value={plant.ps_location} />
                <DetailRow label="Type" value={PLANT_TYPES[plant.ps_type] || 'Unknown'} />
                <DetailRow label="Status" value={plant.ps_fault_status === 3 ? 'Normal' : 'Fault'} isStatus />
                <DetailRow label="Online" value={plant.online_status === 1 ? 'Online' : 'Offline'} />
                <DetailRow label="Installed" value={plant.install_date?.split(' ')[0] || '-'} />
            </div>

            <PlantDeviceList ps_id={plant.ps_id} />
        </div>
    )
}

function DetailRow({
    label,
    value,
    isStatus
}: {
    label: string
    value: string
    isStatus?: boolean
}) {
    return (
        <div className="detail-row">
            <span className="detail-label">{label}</span>
            <span className={`detail-value ${isStatus ? (value === 'Normal' ? 'status-ok' : 'status-err') : ''}`}>
                {value}
            </span>
        </div>
    )
}
