/* export interface Field {
    // ... field properties
} */

/* export class FieldImpl implements Field {
    // ... field implementation
} */

export class Collection {
    name: string;
    num_documents: number;
 //   fields: Field[];
    
    static createFrom(source: Partial<Collection> = {}): Collection {
        return new Collection(source);
    }
    
    constructor(source: Record<string, unknown> = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.name = source["name"] as string;
        this.num_documents = source["num_documents"] as number;
   //     this.fields = this.convertValues(source["fields"], Field);
    }
    
 /*    convertValues<T>(a: unknown, classs: new () => T, asMap = false): T | Record<string, T> {
        if (Array.isArray(a)) {
            return a.map(elem => this.convertValues(elem, classs, asMap)) as T;
        } else if (typeof a === "object") {
            if (asMap) {
                for (const key of Object.keys(a)) {
                    a[key] = new classs(a[key]);
                }
                return a;
            }
            return new classs(a);
        }
        return a;
    } */
} 