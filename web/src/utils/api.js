import { API } from '@/api'
import { vm } from '@/main'
import router from '@/router'
import jwt_decode from "jwt-decode"
import { join_queries, object_to_query } from '@/utils/query'

export function get_data(endpoint, query = null, options = {}, callback = null, callback_arguments = null) {
  var query_str = null
  var url = `/${endpoint}`
  if (query) {
    query_str = object_to_query({s: JSON.stringify(query), ...options})
    url = `/${endpoint}?${query_str}`
  }
  console.log(`GET ${url}`)
  return API
    .get(url)
    .then(response => {
      console.log(response)
      if (callback) {
        return callback(response, callback_arguments)
      }
      if (response.data == undefined) {
        show_feedback(response)
        throw `No data found at /${endpoint}`
      }
    })
    .catch(error => {
      console.log(error)
      throw error
    })
}

// Submit data to an endpoint
export function submit(endpoint, data) {
  var filtered_object = Object.assign({}, data)
  Object.keys(filtered_object).forEach((key, ) => {
    if (key[0] == '_') {
      delete filtered_object[key]
    }
  })
  console.log(`POST /${endpoint}`)
  API
    .post(`/${endpoint}`, [filtered_object])
    .then(response => {
      console.log(response)
      vm.text_alert(`Updated object ${response.data["uid"]}`)
    })
    .catch(error => {
      console.log(error)
    })
}

export function preprocess_data(data) {
  var filtered_object = Object.assign({}, data)
  Object.keys(filtered_object).forEach((key, ) => {
    if (key[0] == '_') {
      delete filtered_object[key]
    }
  })
  return filtered_object
}

export function show_feedback(response, title = null) {
  if(response.data) {
    if (title) {
      vm.text_alert(`Succeeded to ${title} ${response.data.count || ((response.data.data.replaced || []).length + (response.data.data.updated || []).length + (response.data.data.added || []).length)} object(s)`, 'success', title + ' success')
    } else {
      vm.text_alert('Operation successful', 'success')
    }
  } else {
    var message = ''
    if(response.response && response.response.data.description) {
      message = response.response.data.description
    } else {
      if (title) {
        message = `Failed to ${title} object(s): ${response.statusText}`
      } else {
        message = 'An error occurred'
      }
    }
    if (title) {
      vm.text_alert(message, 'danger', title + ' failure')
    } else {
      vm.text_alert(message, 'danger')
    }
  }
}

export function add_items(endpoint, items, callback = null, callback_arguments = null) {
  items = items.map(item => preprocess_data(item))
  console.log(`POST ${endpoint}`)
  console.log(items)
  API
    .post(`/${endpoint}`, items)
    .then(response => {
      console.log(response)
      if (callback) {
        callback(response, callback_arguments)
      }
      show_feedback(response, 'Add')
    })
    .catch(error => console.log(error))
}

export function update_items(endpoint, items, callback = null, callback_arguments = null) {
  items = items.map(item => preprocess_data(item))
  console.log(`PUT ${endpoint}`)
  console.log(items)
  API
    .put(`/${endpoint}`, items)
    .then(response => {
      console.log(response)
      if (callback) {
        callback(response, callback_arguments)
      }
      show_feedback(response, 'Update')
    })
    .catch(error => console.log(error))
}

export function delete_items(endpoint, items, callback = null, callback_arguments = null) {
  var uids = items.map(x => x["uid"])
  uids.forEach(uid => {
    console.log(`DELETE ${uid}`)
  })
  var queries = uids.map(uid => ["=", "uid", uid])
  var query = {
    s: JSON.stringify(join_queries(queries, "OR")),
  }
  var query_str = object_to_query(query)
  API
    .delete(`/${endpoint}?${query_str}`)
    .then(response => {
      console.log(response)
      if (callback) {
        callback(response, callback_arguments)
      }
      show_feedback(response, 'Delete')
    })
    .catch(error => {
      console.log(error)
    })
}

export function to_clipboard(txt) {
  var textArea = document.createElement("textarea");
  textArea.style.position = 'fixed';
  textArea.style.top = 0;
  textArea.style.left = 0;
  textArea.style.width = '2em';
  textArea.style.height = '2em';
  textArea.style.padding = 0;
  textArea.style.border = 'none';
  textArea.style.outline = 'none';
  textArea.style.boxShadow = 'none';
  textArea.style.background = 'transparent';
  textArea.value = txt;

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    document.execCommand('copy')
  } catch (err) {
    console.log('Unable to copy');
  }
  document.body.removeChild(textArea);
}

export function get_alert_color(type) {
  switch (type) {
    case 'ack':
      return 'success'
    case 'esc':
      return 'warning'
    case 'close':
      return 'tertiary'
    case 'open':
      return 'quaternary'
    default:
      return 'primary'
  }
}

export function get_alert_icon(type) {
  switch (type) {
    case 'ack':
      return 'la-thumbs-up'
    case 'esc':
      return 'la-exclamation'
    case 'close':
      return 'la-lock'
    case 'open':
      return 'la-lock-open'
    default:
      return 'la-comment-dots'
  }
}

export function get_alert_tooltip(type) {
  switch (type) {
    case 'ack':
      return 'Acknowledge'
    case 'esc':
      return 'Re-escalate'
    case 'close':
      return 'Close'
    case 'open':
      return 'Re-open'
    default:
      return 'Comment'
  }
}

export function safe_jwt_decode(token, redirect = true) {
  var decoded_token = ''
  try {
    decoded_token = jwt_decode(token)
    return decoded_token
  } catch (error) {
    if (redirect) {
      var return_to = encodeURIComponent(router.currentRoute.value.fullPath)
      router.push('/login?return_to='+return_to)
    } else {
      return ''
    }
  }
}

