import { Route, IndexRoute } from 'react-router'

import CoreLayout from 'layouts/CoreLayout'
import AdminView from 'views/AdminView'
import Master from 'views/master';

export default (
  <Route path='/' component={CoreLayout}>
    <IndexRoute component={AdminView} />
    <Route path='/admin' component={AdminView} />
  </Route>
)
