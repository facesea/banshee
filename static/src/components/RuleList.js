import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'

import Paper from 'material-ui/lib/paper'
import Toolbar from 'material-ui/lib/toolbar/toolbar'
import ToolbarGroup from 'material-ui/lib/toolbar/toolbar-group'
import ToolbarSeparator from 'material-ui/lib/toolbar/toolbar-separator'
import ToolbarTitle from 'material-ui/lib/toolbar/toolbar-title'
import RaisedButton from 'material-ui/lib/raised-button'
import IconButton from 'material-ui/lib/icon-button'
import List from 'material-ui/lib/lists/list'
import ListItem from 'material-ui/lib/lists/list-item'
import FlatButton from 'material-ui/lib/flat-button'
import TextField from 'material-ui/lib/text-field'
import Dialog from 'material-ui/lib/dialog'
import Checkbox from 'material-ui/lib/checkbox'
import Snackbar from 'material-ui/lib/snackbar'

import ActionDelete from 'material-ui/lib/svg-icons/action/delete'

import { actions as projectDetailActions } from '../redux/modules/projectDetail'

import ruleStyles from '../styles/rule.scss'

const mapStateToProps = (state) => ({
  project: state.projectDetail.project,
  rules: state.projectDetail.rules,
  id: state.router.path.split('/')[1],
  ruleOpen: state.projectDetail.ruleOpen,
  opt: state.projectDetail.opt,
  submitDisabled: state.projectDetail.submitDisabled,
  patternErrorText: state.projectDetail.patternErrorText,
  onTrendUp: state.projectDetail.onTrendUp,
  onTrendDown: state.projectDetail.onTrendDown,
  onValueGt: state.projectDetail.onValueGt,
  onValueLt: state.projectDetail.onValueLt,
  onTrendUpAndValueGt: state.projectDetail.onTrendUpAndValueGt,
  onTrendDownAndValueLt: state.projectDetail.onTrendDownAndValueLt,
  thresholdMax: state.projectDetail.thresholdMax,
  thresholdMin: state.projectDetail.thresholdMin,
  trustline: state.projectDetail.trustline,
  thresholdMaxErrorText: state.projectDetail.thresholdMaxErrorText,
  thresholdMinErrorText: state.projectDetail.thresholdMinErrorText,
  trustlineErrorText: state.projectDetail.trustlineErrorText,

  snackbarOpen: state.projectDetail.snackbarOpen,
  snackbarMessage: state.projectDetail.snackbarMessage,

  formField: state.projectDetail.formField,
  formFieldErrorText: state.projectDetail.formFieldErrorText
})

export class RuleList extends React.Component {
  static propTypes = {
    rules: React.PropTypes.array.isRequired,
    getProjectById: React.PropTypes.func.isRequired,
    getAllRules: React.PropTypes.func.isRequired,
    ruleDialogClose: React.PropTypes.func.isRequired,
    id: React.PropTypes.string.isRequired,
    ruleOpen: React.PropTypes.bool.isRequired,
    addRule: React.PropTypes.func.isRequired,
    project: React.PropTypes.object.isRequired,
    opt: React.PropTypes.string.isRequired,
    snackbarOpen: React.PropTypes.bool.isRequired,
    snackbarMessage: React.PropTypes.string.isRequired,
    formField: React.PropTypes.string.isRequired,

    handlePatternChange: React.PropTypes.func.isRequired,
    handleCheck: React.PropTypes.func.isRequired,
    handleInput: React.PropTypes.func.isRequired,
    handleSnackbarClose: React.PropTypes.func.isRequired,
    ruleDialogOpen: React.PropTypes.func.isRequired,
    editProjectName: React.PropTypes.func.isRequired,
    deleteRule: React.PropTypes.func.isRequired,

    submitDisabled: React.PropTypes.bool.isRequired,
    onTrendUp: React.PropTypes.bool.isRequired,
    onTrendDown: React.PropTypes.bool.isRequired,
    onValueGt: React.PropTypes.bool.isRequired,
    onValueLt: React.PropTypes.bool.isRequired,
    onTrendUpAndValueGt: React.PropTypes.bool.isRequired,
    onTrendDownAndValueLt: React.PropTypes.bool.isRequired,
    thresholdMax: React.PropTypes.number,
    thresholdMin: React.PropTypes.number,
    trustline: React.PropTypes.number,

    patternErrorText: React.PropTypes.string.isRequired,
    thresholdMaxErrorText: React.PropTypes.string.isRequired,
    thresholdMinErrorText: React.PropTypes.string.isRequired,
    trustlineErrorText: React.PropTypes.string.isRequired,
    formFieldErrorText: React.PropTypes.string.isRequired
  }

  componentDidMount () {
    let id = this.props.id
    this.props.getProjectById(id)
    this.props.getAllRules(id)
  }

  render () {
    const styles = {
      leftBtn: {
        marginRight: 0
      },
      toolbar: {
        backgroundColor: '#fff'
      },
      separator: {
        margin: '0 16px 0 0'
      },
      link: {
        textDecoration: 'none',
        color: '#444'
      },
      hover: {
        textDecoration: 'underline'
      }
    }

    let ruleActions = [
      <FlatButton
        label='Cancel'
        secondary
        onTouchTap={this.props.ruleDialogClose} />
    ]

    const titles = {
      add: 'Add Rule',
      edit: 'Edit Name'
    }

    let item

    if (this.props.opt === 'add') {
      item = <form id='form' onSubmit={this.props.addRule}>
          <div className={ruleStyles.row}>
            <label className={ruleStyles.label}>Pattern:</label>
            <div className={ruleStyles.rightPart}>
              <TextField
                className={ruleStyles.verticalAlign}
                hintText='timer.count_ps.*'
                onChange={this.props.handlePatternChange}
                errorText={this.props.patternErrorText}/>
            </div>
          </div>

          <div className={ruleStyles.row}>
            <label className={ruleStyles.label}>Alerting:</label>
            <div className={ruleStyles.rightPart}>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onTrendUp', checked) }}
                label='On Trend up'/>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onTrendDown', checked) }}
                label='On Trend down'/>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onValueGt', checked) }}
                label='On value >= thresholdMax'/>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onValueLt', checked) }}
                label='On value <= thresholdMin'/>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onTrendUpAndValueGt', checked) }}
                label='On Trend up and value >= thresholdMax'/>
              <Checkbox
                onCheck={(e, checked) => { this.props.handleCheck('onTrendDownAndValueLt', checked) }}
                label='On Trend down and value <= thresholdMin'/>
              <div>
                <label>thresholdMax:</label>
                <TextField
                  value={this.props.thresholdMax}
                  className={ruleStyles.smallField}
                  onChange={(e) => { this.props.handleInput('thresholdMax', e.target.value) }}
                  errorText={this.props.thresholdMaxErrorText}/>
              </div>
              <div>
                <label>thresholdMin:</label>
                <TextField
                  value={this.props.thresholdMin}
                  className={ruleStyles.smallField}
                  onChange={(e) => { this.props.handleInput('thresholdMin', e.target.value) }}
                  errorText={this.props.thresholdMinErrorText}/>
              </div>
            </div>
          </div>

          <div className={ruleStyles.row}>
            <label className={ruleStyles.label}>Trustline:</label>
            <div className={ruleStyles.rightPart && ruleStyles.divVertical}>
              <label>Don't alert me when value is less than</label>
              <TextField
                  value={this.props.trustline}
                  className={ruleStyles.smallField}
                  onChange={(e) => { this.props.handleInput('trustline', e.target.value) }}
                  errorText={this.props.trustlineErrorText}/>
            </div>
          </div>

        </form>

      ruleActions = [
        ...ruleActions,
        <FlatButton
            label='Submit'
            primary
            form='form'
            disabled={this.props.submitDisabled}
            onTouchTap={this.props.addRule} />
      ]
    }

    if (this.props.opt === 'edit') {
      item = <form id='form' onSubmit={this.props.editProjectName}>
        <TextField
          value={this.props.formField}
          onChange={(e) => { this.props.handleInput('formField', e.target.value) }}
          errorText={this.props.formFieldErrorText}/>
      </form>
      ruleActions = [
        ...ruleActions,
        <FlatButton
            label='Submit'
            primary
            form='form'
            disabled={this.props.submitDisabled}
            onTouchTap={this.props.editProjectName} />
      ]
    }
    console.log(this.props)
    return (
      <Paper zDepth={1}>
        <Toolbar style={styles.toolbar}>
          <ToolbarGroup float='left'>
            <ToolbarTitle text={this.props.project.name} />
            <ToolbarSeparator style={styles.separator}/>
            <ToolbarTitle text='Rules' />
          </ToolbarGroup>
          <ToolbarGroup float='right'>
            <RaisedButton label='Edit Name' primary style={styles.leftBtn} onTouchTap={() => { this.props.ruleDialogOpen('edit') }}/>
            <ToolbarSeparator />
            <RaisedButton label='Add Rule' primary onTouchTap={() => { this.props.ruleDialogOpen('add') }}/>
          </ToolbarGroup>
        </Toolbar>
        <List>
          {
            this.props.rules.map((el, index) => {
              return <ListItem
                primaryText={<Link className={ruleStyles.link} to='/'>{el.pattern}</Link>}
                secondaryText={el.repr}
                rightIconButton={<IconButton tooltip='delete' tooltipPosition='top-right' onClick={() => { this.props.deleteRule(el.id, index) }}> <ActionDelete /> </IconButton>}
                key={el.id}/>
            })
          }

        </List>
        <Dialog
          title={titles[this.props.opt]}
          actions={ruleActions}
          modal={false}
          open={this.props.ruleOpen}
          onRequestClose={this.props.ruleDialogClose}>
          { item }
        </Dialog>

        <Snackbar
          open={this.props.snackbarOpen}
          message={this.props.snackbarMessage}
          autoHideDuration={2000}
          action='ERROR'
          onRequestClose={this.props.handleSnackbarClose}
        />
      </Paper>
    )
  }
}

// export default RuleList
export default connect(mapStateToProps, projectDetailActions)(RuleList)
