export const checksums = {
  "docs": "v3.5.0--quFkNIUZZFAwcn0ok74-KsIERem9u0p5DW-cqEgxrPA"
}
export const checksumsStructure = {
  "docs": "quFkNIUZZFAwcn0ok74-KsIERem9u0p5DW-cqEgxrPA"
}

export const tables = {
  "docs": "_content_docs",
  "info": "_content_info"
}

export default {
  "docs": {
    "type": "page",
    "fields": {
      "id": "string",
      "title": "string",
      "body": "json",
      "description": "string",
      "extension": "string",
      "meta": "json",
      "navigation": "json",
      "path": "string",
      "seo": "json",
      "stem": "string"
    }
  },
  "info": {
    "type": "data",
    "fields": {}
  }
}