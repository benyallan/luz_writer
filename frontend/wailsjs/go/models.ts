export namespace model {
	
	export class Problem {
	    severity: string;
	    code: string;
	    message: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new Problem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.severity = source["severity"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.source = source["source"];
	    }
	}
	export class BuildResult {
	    success: boolean;
	    outputPath: string;
	    problems: Problem[];
	    logTail: string;
	
	    static createFrom(source: any = {}) {
	        return new BuildResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.outputPath = source["outputPath"];
	        this.problems = this.convertValues(source["problems"], Problem);
	        this.logTail = source["logTail"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChapterMeta {
	    id: string;
	    title: string;
	    role: string;
	    wordCount: number;
	    hasOverrides: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ChapterMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.role = source["role"];
	        this.wordCount = source["wordCount"];
	        this.hasOverrides = source["hasOverrides"];
	    }
	}
	export class CustomStyle {
	    id: string;
	    name: string;
	    italic: boolean;
	    bold: boolean;
	    smallCaps: boolean;
	    color?: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomStyle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.italic = source["italic"];
	        this.bold = source["bold"];
	        this.smallCaps = source["smallCaps"];
	        this.color = source["color"];
	    }
	}
	export class FieldOption {
	    value: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new FieldOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.label = source["label"];
	    }
	}
	export class FormField {
	    key: string;
	    label: string;
	    type: string;
	    default?: any;
	    options?: FieldOption[];
	
	    static createFrom(source: any = {}) {
	        return new FormField(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.default = source["default"];
	        this.options = this.convertValues(source["options"], FieldOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FormSchema {
	    fields: FormField[];
	
	    static createFrom(source: any = {}) {
	        return new FormSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fields = this.convertValues(source["fields"], FormField);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PluginManifest {
	    name: string;
	    displayName: string;
	    description: string;
	    core: boolean;
	    documentScope: boolean;
	    schema: FormSchema;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PluginManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.description = source["description"];
	        this.core = source["core"];
	        this.documentScope = source["documentScope"];
	        this.schema = this.convertValues(source["schema"], FormSchema);
	        this.enabled = source["enabled"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Variable {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Variable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class Project {
	    luzVersion: number;
	    title: string;
	    subtitle: string;
	    authors: string[];
	    language: string;
	    chapterOrder: string[];
	    activeTarget: string;
	    variables: Variable[];
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.luzVersion = source["luzVersion"];
	        this.title = source["title"];
	        this.subtitle = source["subtitle"];
	        this.authors = source["authors"];
	        this.language = source["language"];
	        this.chapterOrder = source["chapterOrder"];
	        this.activeTarget = source["activeTarget"];
	        this.variables = this.convertValues(source["variables"], Variable);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Target {
	    luzVersion: number;
	    id: string;
	    name: string;
	    kind: string;
	    documentClass: string;
	    fontSize: string;
	    includeToc: boolean;
	    pluginConfig: Record<string, Array<number>>;
	
	    static createFrom(source: any = {}) {
	        return new Target(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.luzVersion = source["luzVersion"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.kind = source["kind"];
	        this.documentClass = source["documentClass"];
	        this.fontSize = source["fontSize"];
	        this.includeToc = source["includeToc"];
	        this.pluginConfig = source["pluginConfig"];
	    }
	}
	
	export class WorkspaceInfo {
	    path: string;
	    project: Project;
	
	    static createFrom(source: any = {}) {
	        return new WorkspaceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.project = this.convertValues(source["project"], Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

