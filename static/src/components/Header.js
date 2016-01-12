import { Link } from 'react-router'
import AppBar from 'material-ui/lib/app-bar'
import styles from '../styles/Header.scss'

export default class Header extends React.Component {
  render () {
    return (
      <AppBar
        id='header'
        iconElementLeft={<Link to='/' className={styles.brand}>Banshee</Link>}
        iconElementRight={<div className={styles['header-right']}><span className='margin-right-10'>Version2.0.3</span> <Link to='/' className={styles['header-back']}>index</Link> </div>}
      />
    )
  }
}
