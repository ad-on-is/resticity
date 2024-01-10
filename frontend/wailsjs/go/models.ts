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

}

