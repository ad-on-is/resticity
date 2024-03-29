export namespace internal {
	
	export interface AppSettingsNotifications {
	    on_schedule_error: boolean;
	    on_schedule_success: boolean;
	    on_schedule_start: boolean;
	}
	export interface AppSettingsHooks {
	    on_schedule_error: string;
	    on_schedule_success: string;
	    on_schedule_start: string;
	}
	export interface AppSettings {
	    theme: string;
	    preserve_error_logs_days: number;
	    hooks: AppSettingsHooks;
	    notifications: AppSettingsNotifications;
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
	    last_run: string;
	    last_error: string;
	}
	export interface Options {
	    s3_key: string;
	    s3_secret: string;
	    azure_account_name: string;
	    azure_account_key: string;
	    azure_account_sas: string;
	    google_project_id: string;
	    google_application_credentials: string;
	}
	export interface Repository {
	    id: string;
	    name: string;
	    type: string;
	    prune_params: string[][];
	    path: string;
	    password: string;
	    password_file: string;
	    // Go type: Options
	    options: any;
	}
	export interface Config {
	    repositories: Repository[];
	    backups: Backup[];
	    schedules: Schedule[];
	    app_settings: AppSettings;
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
	
	
	export interface ScheduleObject {
	    schedule: Schedule;
	    to_repository?: Repository;
	    from_repository?: Repository;
	    backup?: Backup;
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

