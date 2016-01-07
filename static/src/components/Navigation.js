import List from 'material-ui/lib/lists/list';
import ListItem from 'material-ui/lib/lists/list-item';
import FileFolder from 'material-ui/lib/svg-icons/file/folder'
import CmtContactMail from 'material-ui/lib/svg-icons/communication/contact-mail'
import ActionBuild from 'material-ui/lib/svg-icons/action/build'

import MobileTearSheet from './MobileTearSheet';

export default class Navigation extends React.Component {

  constructor(props) {
    super(props)
    this.state = {open: true}
  }

  render () {
    return (
      <MobileTearSheet>
        <List subheader="Navigation">
          <ListItem primaryText="Project" leftIcon={<FileFolder />} />
          <ListItem primaryText="Receiver" leftIcon={<CmtContactMail />} />
          <ListItem primaryText="Configuration" leftIcon={<ActionBuild />} />
        </List>
      </MobileTearSheet>
    )
  }
}
