import List from 'material-ui/lib/lists/list'
import ListItem from 'material-ui/lib/lists/list-item'
import FileFolder from 'material-ui/lib/svg-icons/file/folder'
import CmtContactMail from 'material-ui/lib/svg-icons/communication/contact-mail'
import ActionBuild from 'material-ui/lib/svg-icons/action/build'

import {SelectableContainerEnhance} from 'material-ui/lib/hoc/selectable-enhance'
let SelectableList = SelectableContainerEnhance(List)

import MobileTearSheet from './MobileTearSheet'

function wrapState (ComposedComponent) {
  const StateWrapper = React.createClass({
    getInitialState () {
      return {selectedIndex: 1}
    },
    handleUpdateSelectedIndex (e, index) {
      this.setState({
        selectedIndex: index
      })
    },
    render () {
      return (
        <ComposedComponent {...this.props} {...this.state}
          valueLink={{value: this.state.selectedIndex, requestChange: this.handleUpdateSelectedIndex}} />
      )
    }
  })
  return StateWrapper
}

SelectableList = wrapState(SelectableList)

export default class Navigation extends React.Component {

  constructor (props) {
    super(props)
    this.state = {open: true}
  }

  render () {
    return (
      <MobileTearSheet>
        <SelectableList
          value={1}
          subheader='Navigation'>

          <ListItem primaryText='Project' value={1} leftIcon={<FileFolder />} />
          <ListItem primaryText='Receiver' value={2} leftIcon={<CmtContactMail />} />
          <ListItem primaryText='Configuration' value={3} leftIcon={<ActionBuild />} />
        </SelectableList>
      </MobileTearSheet>
    )
  }
}
