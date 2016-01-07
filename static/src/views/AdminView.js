import { Paper } from 'material-ui'
import ClearFix from 'material-ui/lib/clearfix'

import Navigation from '../components/Navigation'
import AdminTable from '../components/AdminTable'

export class AboutView extends React.Component {
  render () {
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
    }

    return (
      <div className='container-admin' style={styles.container}>
        <ClearFix style={styles.clearfix}>
          <Navigation />
        </ClearFix>
        <Paper style={styles.paper}>
          <AdminTable />
        </Paper>

      </div>
    )
  }
}

export default AboutView
