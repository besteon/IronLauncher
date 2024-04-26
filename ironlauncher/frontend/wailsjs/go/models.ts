export namespace main {
	
	export class Settings {
	    romsFolder: string;
	    defaultRom: string;
	    defaultMode: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.romsFolder = source["romsFolder"];
	        this.defaultRom = source["defaultRom"];
	        this.defaultMode = source["defaultMode"];
	    }
	}

}

