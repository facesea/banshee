import { Route, IndexRoute } from 'react-router'

import CoreLayout from 'layouts/CoreLayout'
import AdminView from 'views/AdminView'
import ProjectDetailView from 'views/ProjectDetailView'

export default (
  <Route path='/' component={CoreLayout}>
    <IndexRoute component={AdminView} />
    <Route path='/admin' component={AdminView} />
    <Route path='/admin/:projectId' component={ProjectDetailView} />
  </Route>
)
