export namespace main {
	
	export interface Settings {
	    repositories: restic.Repository[];
	    backups: restic.Backup[];
	}

}

export namespace restic {
	
	export interface Backup {
	    id: string;
	    path: string;
	    name: string;
	    cron: string;
	    backup_params: string[][];
	    targets: string[];
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
	    prune_params: Param[];
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

