export namespace main {

	export class FileNode {
	    name: string;
	    path: string;
	    isDir: boolean;

	    static createFrom(source: any = {}) {
	        return new FileNode(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	    }
	}

	export class RecentProject {
	    name: string;
	    path: string;

	    static createFrom(source: any = {}) {
	        return new RecentProject(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}
