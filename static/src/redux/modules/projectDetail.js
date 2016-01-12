import {
  createAction, handleActions
}
from 'redux-actions'
import request from 'superagent'

// ------------------------------------
// Constants
// ------------------------------------
export const SET_ALL_RULES = 'SET_ALL_RULES'
export const GET_ALL_RULES = 'GET_ALL_RULES'
export const SET_PROJECT = 'SET_PROJECT'
export const RULE_DIALOG_CLOSE = 'RULE_DIALOG_CLOSE'
export const ADD_RULE = 'ADD_RULE'
export const HANDLE_PATTERN_CHANGE = 'HANDLE_PATTERN_CHANGE'
export const HANDLE_CHECK = 'HANDLE_CHECK'
export const HANDLE_INPUT = 'HANDLE_INPUT'
export const ADD_RULE_SUCCESS = 'ADD_RULE_SUCCESS'
export const ADD_RULE_FAIL = 'ADD_RULE_FAIL'
export const RULE_DIALOG_OPEN = 'RULE_DIALOG_OPEN'
export const HANDLE_SNACKBAR_CLOSE = 'HANDLE_SNACKBAR_CLOSE'
export const EDIT_PROJECT_NAME_SUCCESS = 'EDIT_PROJECT_NAME_SUCCESS'
export const DELETE_RULE_SUCCESS = 'DELETE_RULE_SUCCESS'
export const EDIT_PROJECT_NAME_FAIL = 'EDIT_PROJECT_NAME_FAIL'

export const INIT_STATE = {
  project: {},
  rules: [],
  ruleOpen: false,
  opt: '',
  submitDisabled: false,
  onTrendUp: false,
  onTrendDown: false,
  onValueGt: false,
  onValueLt: false,
  onTrendUpAndValueGt: false,
  onTrendDownAndValueLt: false,
  patternErrorText: '',
  thresholdMaxErrorText: '',
  thresholdMinErrorText: '',
  trustlineErrorText: '',
  snackbarOpen: false,
  snackbarMessage: '',

  formField: '',
  formFieldErrorText: ''
}

// ------------------------------------
// Actions
// ------------------------------------
export const setProject = createAction(SET_PROJECT, (project) => project)
export const setAllRules = createAction(SET_ALL_RULES, (rules) => rules)
export const ruleDialogClose = createAction(RULE_DIALOG_CLOSE, () => false)
export const ruleDialogOpen = createAction(RULE_DIALOG_OPEN, (opt) => opt)
export const handlePatternChange = createAction(HANDLE_PATTERN_CHANGE, (e) => e.target.value)
export const handleCheck = createAction(HANDLE_CHECK, (name, checked) => ({name: name, checked: checked}))
export const handleInput = createAction(HANDLE_INPUT, (name, val) => ({name: name, val: val}))
export const addRuleFail = createAction(ADD_RULE_FAIL, (msg) => msg)
export const addRuleSuccess = createAction(ADD_RULE_SUCCESS, (rule) => rule)
export const handleSnackbarClose = createAction(HANDLE_SNACKBAR_CLOSE, () => false)
export const editProjectNameSucess = createAction(EDIT_PROJECT_NAME_SUCCESS, (project) => project)
export const editProjectNameFail = createAction(EDIT_PROJECT_NAME_FAIL, (msg) => msg)
export const deleteRuleSuccess = createAction(DELETE_RULE_SUCCESS, (index) => index)

export const getProjectById = (id) => {
  return (dispatch, getState) => {
    request
      .get('/api/project/' + id)
      .end(function (err, res) {
        if (err || !res.ok) {
          console.error('get project error')
        } else {
          dispatch(setProject(res.body))
        }
      })
  }
}

export const getAllRules = (id) => {
  return (dispatch, getState) => {
    request
      .get('/api/project/' + id + '/rules')
      .end(function (err, res) {
        if (err || !res.ok) {
          console.error('get projects error')
        } else {
          dispatch(setAllRules(res.body))
        }
      })
  }
}

export const editProjectName = (e) => {
  if (e) {
    e.preventDefault()
  }
  return (dispatch, getState) => {
    const state = getState().projectDetail

    request
      .patch('/api/project/' + state.project.id)
      .send({name: state.formField})
      .end(function (err, res) {
        if (err || !res.ok) {
          dispatch(editProjectNameFail(res.body.msg))
        } else {
          dispatch(editProjectNameSucess(res.body))
        }
      })
  }
}

export const addRule = () => {
  return (dispatch, getState) => {
    const state = getState().projectDetail
    const params = {
      'pattern': state.pattern,
      'onTrendUp': state.onTrendUp,
      'onTrendDown': state.onTrendDown,
      'onValueGt': state.onValueGt,
      'onValueLt': state.onValueLt,
      'onTrendUpAndValueGt': state.onTrendUpAndValueGt,
      'onTrendDownAndValueLt': state.onTrendDownAndValueLt,
      'thresholdMax': state.thresholdMax,
      'thresholdMin': state.thresholdMin,
      'trustline': state.trustline
    }
    request
      .post('/api/project/' + state.project.id + '/rule')
      .send(params)
      .end(function (err, res) {
        if (err || !res.ok) {
          dispatch(addRuleFail(err.msg))
        } else {
          dispatch(addRuleSuccess(res.body))
        }
      })
  }
}

export const deleteRule = (id, index) => {
  return (dispatch, getState) => {
    request
      .del('/api/rule/' + id)
      .end(function (err, res) {
        if (err || !res.ok) {
          console.error('delete rule error')
        } else {
          dispatch(deleteRuleSuccess(index))
        }
      })
  }
}

export const actions = {
  getProjectById,
  getAllRules,
  ruleDialogClose,
  addRule,
  editProjectName,
  deleteRule,

  handlePatternChange,
  handleCheck,
  handleInput,
  handleSnackbarClose,
  ruleDialogOpen
}

// ------------------------------------
// Reducer
// ------------------------------------
export default handleActions({
  [SET_PROJECT]: (state, { payload }) => {
    return Object.assign({}, state, {
      project: payload
    })
  },
  [SET_ALL_RULES]: (state, { payload }) => {
    return Object.assign({}, state, {
      rules: payload
    })
  },
  [RULE_DIALOG_CLOSE]: (state, { payload }) => {
    return Object.assign({}, state, {
      ruleOpen: payload
    })
  },
  [HANDLE_PATTERN_CHANGE]: (state, { payload }) => {
    if (!payload) {
      state.patternErrorText = 'This field is required.'
    } else {
      state.patternErrorText = ''
    }

    validate(state)

    return Object.assign({}, state, {
      pattern: payload
    })
  },
  [HANDLE_CHECK]: (state, { payload }) => {
    state[payload.name] = payload.checked
    return Object.assign({}, state)
  },
  [HANDLE_INPUT]: (state, { payload }) => {
    if (state.opt === 'add' && (Number(payload.val) || payload.val === '')) {
      state[payload.name] = Number(payload.val) || undefined
      state[payload.name + 'ErrorText'] = ''
    } else if (state.opt === 'add') {
      state[payload.name] = undefined
      state[payload.name + 'ErrorText'] = 'This field must be numeric.'
    } else if (state.opt === 'edit') {
      state[payload.name] = payload.val
      state[payload.name + 'ErrorText'] = payload.val ? '' : 'This field is required.'
    }

    validate(state)

    return Object.assign({}, state)
  },
  [ADD_RULE_SUCCESS]: (state, { payload }) => {
    return Object.assign({}, state, {
      rules: [
        ...state.rules,
        payload
      ],
      ruleOpen: false,
      submitDisabled: false,
      onTrendUp: false,
      onTrendDown: false,
      onValueGt: false,
      onValueLt: false,
      onTrendUpAndValueGt: false,
      onTrendDownAndValueLt: false,
      patternErrorText: '',
      thresholdMaxErrorText: '',
      thresholdMinErrorText: '',
      trustlineErrorText: ''
    })
  },
  [ADD_RULE_FAIL]: (state, { payload }) => {
    return Object.assign({}, state, {
      snackbarOpen: true,
      snackbarMessage: payload
    })
  },
  [RULE_DIALOG_OPEN]: (state, { payload }) => {
    if (payload === 'edit') {
      state.formField = state.project.name
      state.formFieldErrorText = ''
    }

    return Object.assign({}, state, {
      ruleOpen: true,
      opt: payload
    })
  },
  [HANDLE_SNACKBAR_CLOSE]: (state, { payload }) => {
    return Object.assign({}, state, {
      snackbarOpen: payload
    })
  },
  [EDIT_PROJECT_NAME_SUCCESS]: (state, { payload }) => {
    return Object.assign({}, state, {
      ruleOpen: false,
      project: payload
    })
  },
  [EDIT_PROJECT_NAME_FAIL]: (state, { payload }) => {
    return Object.assign({}, state, {
      snackbarOpen: true,
      snackbarMessage: payload
    })
  },
  [DELETE_RULE_SUCCESS]: (state, { payload }) => {
    state.rules = state.rules.filter((el, index) => {
      return index !== payload
    })

    return Object.assign({}, state)
  }
},

// init state
INIT_STATE
)

// validate form
function validate (state) {
  if (state.patternErrorText || state.thresholdMaxErrorText || state.thresholdMinErrorText || state.trustlineErrorText || state.formFieldErrorText) {
    state.submitDisabled = true
  } else {
    state.submitDisabled = false
  }
}
