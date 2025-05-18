export namespace core {
	
	export class ImageFIle {
	    Code: number;
	    Name: string;
	    Id: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageFIle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.Name = source["Name"];
	        this.Id = source["Id"];
	    }
	}
	export class Message {
	    Code: number;
	    Msg: string;
	    Data: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.Msg = source["Msg"];
	        this.Data = source["Data"];
	    }
	}

}

