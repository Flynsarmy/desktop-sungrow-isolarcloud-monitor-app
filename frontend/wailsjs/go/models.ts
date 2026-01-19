export namespace main {
	
	export class Credentials {
	    appKey: string;
	    secretKey: string;
	    authUrl: string;
	    accessToken?: string;
	    refreshToken?: string;
	    tokenExpiry?: number;
	    gatewayUrl?: string;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appKey = source["appKey"];
	        this.secretKey = source["secretKey"];
	        this.authUrl = source["authUrl"];
	        this.accessToken = source["accessToken"];
	        this.refreshToken = source["refreshToken"];
	        this.tokenExpiry = source["tokenExpiry"];
	        this.gatewayUrl = source["gatewayUrl"];
	    }
	}
	export class Plant {
	    ps_id: number;
	    ps_name: string;
	    description?: string;
	    ps_type: number;
	    online_status: number;
	    valid_flag: number;
	    grid_connection_status: number;
	    install_date: string;
	    ps_location: string;
	    latitude: number;
	    longitude: number;
	    ps_fault_status: number;
	    connect_type: number;
	    update_time: string;
	    ps_current_time_zone: string;
	    grid_connection_time?: string;
	    build_status: number;
	    today_energy?: string;
	
	    static createFrom(source: any = {}) {
	        return new Plant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ps_id = source["ps_id"];
	        this.ps_name = source["ps_name"];
	        this.description = source["description"];
	        this.ps_type = source["ps_type"];
	        this.online_status = source["online_status"];
	        this.valid_flag = source["valid_flag"];
	        this.grid_connection_status = source["grid_connection_status"];
	        this.install_date = source["install_date"];
	        this.ps_location = source["ps_location"];
	        this.latitude = source["latitude"];
	        this.longitude = source["longitude"];
	        this.ps_fault_status = source["ps_fault_status"];
	        this.connect_type = source["connect_type"];
	        this.update_time = source["update_time"];
	        this.ps_current_time_zone = source["ps_current_time_zone"];
	        this.grid_connection_time = source["grid_connection_time"];
	        this.build_status = source["build_status"];
	        this.today_energy = source["today_energy"];
	    }
	}
	export class PlantDevice {
	    uuid: number;
	    ps_key: string;
	    device_sn: string;
	    device_name: string;
	    device_type: number;
	    type_name: string;
	    device_model_id: number;
	    device_model_code: string;
	    dev_fault_status: number;
	    dev_status: string;
	    claim_state: number;
	    device_code: number;
	    chnnl_id: number;
	    communication_dev_sn: string;
	    ps_id: number;
	
	    static createFrom(source: any = {}) {
	        return new PlantDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.ps_key = source["ps_key"];
	        this.device_sn = source["device_sn"];
	        this.device_name = source["device_name"];
	        this.device_type = source["device_type"];
	        this.type_name = source["type_name"];
	        this.device_model_id = source["device_model_id"];
	        this.device_model_code = source["device_model_code"];
	        this.dev_fault_status = source["dev_fault_status"];
	        this.dev_status = source["dev_status"];
	        this.claim_state = source["claim_state"];
	        this.device_code = source["device_code"];
	        this.chnnl_id = source["chnnl_id"];
	        this.communication_dev_sn = source["communication_dev_sn"];
	        this.ps_id = source["ps_id"];
	    }
	}

}

