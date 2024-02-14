export namespace main {
	
	export enum RepositoryType {
	    LOCAL = "Local",
	    S3 = "S3",
	    AZURE = "Azure",
	    GOOGLE = "Google",
	}
	export interface Backup {
	    id: string;
	    path: string;
	    name: string;
	    cron: string;
	    backup_params: string[][];
	    targets: string[];
	}
	export interface Schedule {
	    id: string;
	    action: string;
	    backup_id: string;
	    to_repository_id: string;
	    from_repository_id: string;
	    cron: string;
	    active: boolean;
	}
	export interface Options {
	    s3_key: string;
	    s3_secret: string;
	    azure_account_name: string;
	    azure_account_key: string;
	    azure_account_sas: string;
	    azure_endpoint_suffix: string;
	}
	export interface Repository {
	    id: string;
	    name: string;
	    type: RepositoryType;
	    prune_params: string[][];
	    path: string;
	    password: string;
	    // Go type: Options
	    options: any;
	}
	export interface Mount {
	    repository_id: string;
	    path: string;
	}
	export interface Config {
	    mounts: Mount[];
	    repositories: Repository[];
	    backups: Backup[];
	    schedules: Schedule[];
	}
	export interface FileDescriptor {
	    name: string;
	    type: string;
	    path: string;
	    size: number;
	    mtime: string;
	}
	export interface GroupKey {
	    hostname: string;
	    paths: string[];
	    tags: string[];
	}
	
	
	
	export interface Snapshot {
	    id: string;
	    // Go type: time
	    time: any;
	    paths: string[];
	    hostname: string;
	    username: string;
	    uid: number;
	    gid: number;
	    short_id: string;
	    tags: string[];
	    program_version: string;
	}
	export interface SnapshotGroup {
	    group_key: GroupKey;
	    snapshots: Snapshot[];
	}

}

