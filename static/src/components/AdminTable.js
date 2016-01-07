import React from 'react'
import {
  Table,
  TableBody,
  TableHeader,
  TableFooter,
  TableRow,
  TableHeaderColumn,
  TableRowColumn,
  TextField,
  Toggle,
} from 'material-ui'

export class AdminTable extends React.Component {
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
      }
    };

    return (
      <Table
        height={'300px'}
        fixedHeader={true}
        fixedFooter={true}
        selectable={true}
        multiSelectable={false}
        >
        <TableHeader enableSelectAll={true}>
          <TableRow>
            <TableHeaderColumn colSpan="3" tooltip="Super Header" style={{textAlign: 'center'}}>
              Projects
            </TableHeaderColumn>
          </TableRow>
          <TableRow>
            <TableHeaderColumn tooltip="ID">ID</TableHeaderColumn>
            <TableHeaderColumn tooltip="Project Name">Name</TableHeaderColumn>
            <TableHeaderColumn tooltip="opt">Opt</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody
          deselectOnClickaway={false}
          showRowHover={true}
          stripedRows={true}>
          <TableRow selected={true}>
            <TableRowColumn>1</TableRowColumn>
            <TableRowColumn>John Smith</TableRowColumn>
            <TableRowColumn>del</TableRowColumn>
          </TableRow>
        </TableBody>
      </Table>
    )
  }
}

export default AdminTable
