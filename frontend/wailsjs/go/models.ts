export namespace main {
	
	export class Signals {
	    index: number;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new Signals(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.label = source["label"];
	    }
	}

}

