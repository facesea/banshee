import React from 'react'
import {
  Table,
  TableBody,
  TableHeader,
  TableRow,
  TableHeaderColumn,
  TableRowColumn
} from 'material-ui'

export class AdminTable extends React.Component {
  render () {
    return (
      <Table
        height={'300px'}
        fixedHeader
        fixedFooter
        selectable
        >
        <TableHeader enableSelectAll>
          <TableRow>
            <TableHeaderColumn colSpan='3' tooltip='Super Header' style={{textAlign: 'center'}}>
              Projects
            </TableHeaderColumn>
          </TableRow>
          <TableRow>
            <TableHeaderColumn tooltip='ID'>ID</TableHeaderColumn>
            <TableHeaderColumn tooltip='Project Name'>Name</TableHeaderColumn>
            <TableHeaderColumn tooltip='opt'>Opt</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody
          showRowHover
          stripedRows>
          <TableRow selected>
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
