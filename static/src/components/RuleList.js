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

import ActionDelete from 'material-ui/lib/svg-icons/action/delete'

import { actions as projectDetailActions } from '../redux/modules/projectDetail'

import ruleStyles from '../styles/rule.scss'

const mapStateToProps = (state) => ({
  project: state.projectDetail.project,
  rules: state.projectDetail.rules,
  id: state.router.path.split('/')[1]
})

export class RuleList extends React.Component {
  static propTypes = {
    getProjectById: React.PropTypes.func.isRequired,
    getAllRules: React.PropTypes.func.isRequired,
    id: React.PropTypes.string.isRequired
  };

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

    const rightIconButton = (
      <IconButton tooltip='delete' tooltipPosition='top-right' >
        <ActionDelete />
      </IconButton>
    )

    return (
      <Paper zDepth={1}>
        <Toolbar style={styles.toolbar}>
          <ToolbarGroup float='left'>
            <ToolbarTitle text='Project Detail' />
            <ToolbarSeparator style={styles.separator}/>
            <ToolbarTitle text='Rules' />
          </ToolbarGroup>
          <ToolbarGroup float='right'>
            <RaisedButton label='Edit Name' primary style={styles.leftBtn}/>
            <ToolbarSeparator />
            <RaisedButton label='Add Rule' primary/>
          </ToolbarGroup>
        </Toolbar>
        <List>
          <ListItem
            primaryText={<Link className={ruleStyles.link} to='/'> project name</Link>}
            secondaryText='Change your Google+ profile photo'
            rightIconButton={rightIconButton}/>
        </List>
      </Paper>
    )
  }
}

// export default RuleList
export default connect(mapStateToProps, projectDetailActions)(RuleList)
