import React from 'react'
import { connect } from 'react-redux'
import {
  Table,
  TableBody,
  TableHeader,
  TableRow,
  TableHeaderColumn,
  TableRowColumn,
  RaisedButton
} from 'material-ui'
import Dialog from 'material-ui/lib/dialog'
import FlatButton from 'material-ui/lib/flat-button'
import IconButton from 'material-ui/lib/icon-button'
import TextField from 'material-ui/lib/text-field'
import Snackbar from 'material-ui/lib/snackbar'
import ContentAddBox from 'material-ui/lib/svg-icons/content/add-box'
import { Colors } from 'material-ui/lib/styles'

import { actions as projectActions } from '../redux/modules/project'

const mapStateToProps = (state) => ({
  projects: state.project.projects,
  open: state.project.open,
  projectName: state.project.projectName,
  errorText: state.project.errorText,
  snackbarMessage: state.project.snackbarMessage,
  snackbarOpen: state.project.snackbarOpen
})
export class AdminTable extends React.Component {
  static propTypes = {
    projects: React.PropTypes.array.isRequired,
    open: React.PropTypes.bool.isRequired,
    projectName: React.PropTypes.string.isRequired,
    errorText: React.PropTypes.string.isRequired,
    getAllProjects: React.PropTypes.func.isRequired,
    dialogOpen: React.PropTypes.func.isRequired,
    dialogClose: React.PropTypes.func.isRequired,
    handleInputChange: React.PropTypes.func.isRequired,
    handleSnackbarClose: React.PropTypes.func.isRequired,
    createProject: React.PropTypes.func.isRequired,
    snackbarOpen: React.PropTypes.bool.isRequired,
    snackbarMessage: React.PropTypes.string.isRequired
  }

  componentDidMount () {
    this.props.getAllProjects()
  }

  render () {
    const styles = {
      btnPrimary: {
        minWidth: 50,
        marginRight: 10,
        verticalAlign: 'bottom'
      },
      column: {
        textAlign: 'center'
      },
      iconButton: {
        width: 24,
        height: 24,
        padding: 0,
        float: 'right',
        color: Colors.lightGreen500
      }
    }
    const actions = [
      <FlatButton
        label='Cancel'
        secondary
        onTouchTap={this.props.dialogClose} />,
      <FlatButton
        label='Submit'
        primary
        form='form'
        onTouchTap={this.props.createProject} />
    ]
    return (
      <div>
        <Table
          height={'500px'}
          fixedHeader
          fixedFooter>
          <TableHeader>
            <TableRow>
              <TableHeaderColumn colSpan='3' style={{textAlign: 'center'}}>
                Projects
                <IconButton touch style={styles.iconButton} onTouchTap={this.props.dialogOpen}>
                  <ContentAddBox color={Colors.lightGreen500}/>
                </IconButton>
              </TableHeaderColumn>
            </TableRow>
            <TableRow>
              <TableHeaderColumn>ID</TableHeaderColumn>
              <TableHeaderColumn>Name</TableHeaderColumn>
              <TableHeaderColumn style={styles.column}>Opt</TableHeaderColumn>
            </TableRow>
          </TableHeader>
          <TableBody
            showRowHover
            stripedRows
            selectable={false}
            deselectOnClickaway={false}>
            {
              this.props.projects.map((el, index) => {
                return <TableRow selectable={false} key={el.id}>
                    <TableRowColumn>{el.id}</TableRowColumn>
                    <TableRowColumn>{el.name}</TableRowColumn>
                    <TableRowColumn style={styles.column}>
                      <RaisedButton fullWidth={false} label='View' primary style={styles.btnPrimary} onTouchTap={this.props.getAllProjects}/>
                      <RaisedButton label='Edit' secondary linkButton href={'/admin/' + el.id} style={styles.btnPrimary} />
                    </TableRowColumn>
                  </TableRow>
              })
            }
          </TableBody>
        </Table>
        <Dialog
          title='Create'
          actions={actions}
          modal={false}
          open={this.props.open}
          onRequestClose={this.props.dialogClose}>
            <form id='form' onSubmit={this.props.createProject}>
              <TextField
                hintText='Project Name'
                onChange={this.props.handleInputChange}
                errorText={this.props.errorText}/>
            </form>
        </Dialog>
        <Snackbar
          open={this.props.snackbarOpen}
          message={this.props.snackbarMessage}
          autoHideDuration={2000}
          action='ERROR'
          onRequestClose={this.props.handleSnackbarClose}
        />
      </div>
    )
  }
}

// export default AdminTable
export default connect(mapStateToProps, projectActions)(AdminTable)
