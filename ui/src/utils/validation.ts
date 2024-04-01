/**
 * An object that contain the feedback/status of an alert,
 * and extract it from the response
 */

import { ApiError, ValidationError } from '@/api'

export class ValidationStatus {
  feedback: Map<string, string>
  status: Map<string, string>

  constructor() {
    this.feedback = new Map()
    this.status = new Map()
  }

  onResponse(resp: ApiError) {
    if (resp.status != 422) {
      return
    }
    const body = resp.body as ValidationError
    for (const [field, errors] of Object.entries(body.field_errors)) {
      this.status.set(field, "error")
      this.feedback.set(field, errors.join("\n"))
    }
  }

  reset() {
    this.feedback.clear()
    this.status.clear()
  }
}
