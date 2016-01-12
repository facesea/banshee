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

export const INIT_STATE = {
  project: {},
  rules: [],
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
  trustlineErrorText: '',
  snackbarOpen: false,
  snackbarMessage: ''
}

// ------------------------------------
// Actions
// ------------------------------------
export const setProject = createAction(SET_PROJECT, (project) => project)
export const setAllRules = createAction(SET_ALL_RULES, (rules) => rules)
export const ruleDialogClose = createAction(RULE_DIALOG_CLOSE, () => false)
export const ruleDialogOpen = createAction(RULE_DIALOG_OPEN, () => true)
export const handlePatternChange = createAction(HANDLE_PATTERN_CHANGE, (e) => e.target.value)
export const handleCheck = createAction(HANDLE_CHECK, (name, checked) => ({name: name, checked: checked}))
export const handleInput = createAction(HANDLE_INPUT, (name, val) => ({name: name, val: val}))
export const addRuleFail = createAction(ADD_RULE_FAIL, (msg) => msg)
export const addRuleSuccess = createAction(ADD_RULE_SUCCESS, (rule) => rule)
export const handleSnackbarClose = createAction(HANDLE_SNACKBAR_CLOSE, () => false)

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
          dispatch(addRuleFail(res.body.msg))
        } else {
          dispatch(addRuleSuccess(res.body))
        }
      })
  }
}
export const actions = {
  getProjectById,
  getAllRules,
  ruleDialogClose,
  addRule,

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
    if (Number(payload.val) || payload.val === '') {
      state[payload.name] = Number(payload.val) || undefined
      state[payload.name + 'ErrorText'] = ''
    } else {
      state[payload.name] = undefined
      state[payload.name + 'ErrorText'] = 'This field must be numeric.'
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
    return Object.assign({}, state, {
      ruleOpen: payload
    })
  },
  [HANDLE_SNACKBAR_CLOSE]: (state, { payload }) => {
    return Object.assign({}, state, {
      snackbarOpen: payload
    })
  }
},

// init state
INIT_STATE
)

// validate form
function validate (state) {
  if (state.patternErrorText || state.thresholdMaxErrorText || state.thresholdMinErrorText || state.trustlineErrorText) {
    state.submitDisabled = true
  } else {
    state.submitDisabled = false
  }
}
