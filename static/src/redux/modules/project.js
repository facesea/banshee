import {
  createAction, handleActions
}
from 'redux-actions'
import request from 'superagent'

// ------------------------------------
// Constants
// ------------------------------------
export const SET_ALL_PROJECTS = 'SET_ALL_PROJECTS'
export const DIALOG_OPEN = 'DIALOG_OPEN'
export const DIALOG_CLOSE = 'DIALOG_CLOSE'
export const GET_ALL_PROJECTS_FAIL = 'GET_ALL_PROJECTS_FAIL'
export const CREATE_PROJECT_SUCCESS = 'CREATE_PROJECT_SUCCESS'
export const CREATE_PROJECT_FAIL = 'CREATE_PROJECT_FAIL'
export const HANDLE_INPUT_CHANGE = 'HANDLE_INPUT_CHANGE'
export const HANDLE_SNACKBAR_CLOSE = 'HANDLE_SNACKBAR_CLOSE'

export const INIT_STATE = {
  projects: [],
  projectName: '',
  errorText: '',
  open: false,
  snackbarMessage: 'sss',
  snackbarOpen: false
}

// ------------------------------------
// Actions
// ------------------------------------
export const dialogOpen = createAction(DIALOG_OPEN, () => true)
export const dialogClose = createAction(DIALOG_CLOSE, () => false)
export const setProjects = createAction(SET_ALL_PROJECTS, (projects = []) => projects)
export const getAllProjectsFail = createAction(GET_ALL_PROJECTS_FAIL, (msg) => msg)
export const createProjectSuccess = createAction(CREATE_PROJECT_SUCCESS, (project) => project)
export const createProjectFail = createAction(CREATE_PROJECT_FAIL, (msg) => msg)
export const handleInputChange = createAction(HANDLE_INPUT_CHANGE, (e) => e.target.value)
export const handleSnackbarClose = createAction(HANDLE_SNACKBAR_CLOSE, () => false)

export const getAllProjects = () => {
  return (dispatch, getState) => {
    return request.get('/api/projects')
    .end((err, res) => {
      if (err || !res.ok) {
        dispatch(getAllProjectsFail(res.body.msg))
      } else {
        dispatch(setProjects(res.body))
      }
    })
  }
}

export const createProject = (e) => {
  if (e) {
    e.preventDefault()
  }

  return (dispatch, getState) => {
    let state = getState().project

    return request.post('/api/project')
    .send({
      name: state.projectName
    })
    .end((err, res) => {
      if (err || !res.ok) {
        dispatch(createProjectFail(res.body.msg))
      } else {
        dispatch(createProjectSuccess(res.body))
      }
    })
  }
}

export const actions = {
  setProjects,
  getAllProjects,
  dialogOpen,
  dialogClose,
  createProject,
  createProjectFail,
  handleInputChange,
  handleSnackbarClose
}

// ------------------------------------
// Reducer
// ------------------------------------
export default handleActions({
  [SET_ALL_PROJECTS]: (state, {
    payload
  }) => {
    return Object.assign({}, INIT_STATE, {
      projects: payload
    })
  }, [DIALOG_OPEN]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      open: payload
    })
  }, [DIALOG_CLOSE]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      open: payload
    })
  }, [CREATE_PROJECT_SUCCESS]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      open: false,
      projectName: '',
      projects: [
        ...state.projects,
        {
          id: payload.id,
          name: payload.name
        }
      ]
    })
  }, [HANDLE_INPUT_CHANGE]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      projectName: payload,
      errorText: payload ? '' : 'This field is required.'
    })
  }, [HANDLE_SNACKBAR_CLOSE]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      snackbarOpen: payload
    })
  }, [CREATE_PROJECT_FAIL]: (state, {
    payload
  }) => {
    return Object.assign({}, state, {
      snackbarOpen: true,
      snackbarMessage: payload
    })
  }
},

// init state
INIT_STATE
)
