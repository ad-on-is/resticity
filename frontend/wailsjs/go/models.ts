export namespace main {
	
	export interface Schedule {
	    id: string;
	    backup_id: string;
	    to_repository_id: string;
	    from_repository_id: string;
	    cron: string;
	    active: boolean;
	}
	export interface BackupJob {
	    job_id: number[];
	    schedule: Schedule;
	}
	export interface Options {
	    b2_account_id: string;
	    b2_account_key: string;
	    azure_account_name: string;
	    azure_account_key: string;
	    azure_account_sas: string;
	    azure_endpoint_suffix: string;
	}
	export interface Repository {
	    id: string;
	    name: string;
	    type: number;
	    prune_params: string[][];
	    path: string;
	    password: string;
	    // Go type: Options
	    options: any;
	}
	
	export interface Snapshot {
	    id: string;
	    time: string;
	    paths: string[];
	    hostname: string;
	    username: string;
	    uid: number;
	    gid: number;
	    short_id: string;
	    tags: string[];
	    program_version: string;
	}

}

