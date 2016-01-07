import { Link } from 'react-router'
import { Paper } from 'material-ui'
import ClearFix from 'material-ui/lib/clearfix'

import BellNavigation from '../components/BellNavigation'
import AdminTable from '../components/AdminTable'

export class AboutView extends React.Component {
  render () {
    const navList = [
      {
        name: 'Project',
        link: '/admin/project'
      },
      {
        name: 'Receiver',
        link: '/admin/receiver'
      },
      {
        name: 'Configuration',
        link: '/admin/configuration'
      }
    ];

    const styles = {
      container: {
        padding: 20
      },
      paper: {
        width: 'calc(100% - 225px)',
        float: 'right'
      },
      clearfix: {
        display: 'inline-block'
      }
    };

    return (
      <div className='container-admin' style={styles.container}>
        <ClearFix style={styles.clearfix}>
          <BellNavigation />
        </ClearFix>
        <Paper style={styles.paper}>
          <AdminTable />
        </Paper>

      </div>
    )
  }
}

export default AboutView
