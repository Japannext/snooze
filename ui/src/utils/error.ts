import { v4 as uuidv4 } from 'uuid'
import { ApiError, ValidationError } from '@/api'

export class Alert {
  id: string
  timestamp: Date
  type: string
  title: string
  text: string

  constructor(mytype: string, title: string, text: string) {
    this.id = uuidv4()
    this.timestamp = new Date()
    this.type = mytype
    this.title = `${title} (${this.timestamp.toISOString()})`
    this.text = text
  }
}

export class HttpAlert extends Alert {
  constructor(code: number, title: string, text: string) {
    const mytype = code < 500 ? "warning" : "error"
    super(mytype, `HTTP ${code}: ${title}`, text)
  }

  static fromResponse(error: ApiError): HttpAlert {
    if (error.status == 422) { // ValidationError
      let text = ""
      const body: ValidationError = error.body as ValidationError
      for (const [field, lines] of Object.entries(body.field_errors)) {
        text += `Field '${field}': ` +  lines.join("\n") + "\n"
      }
      return new HttpAlert(error.status, "ValidationError", text)
    } else { // Default
      return new HttpAlert(error.status, error.statusText || error.name, error.body.text)
    }
  }
}

