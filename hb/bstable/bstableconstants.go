package bstable

//Table base class for tables
const Table = "table"

//TableCellStatusActive Applies the hover color to the table row or table cell
const TableCellStatusActive = "active"

//TableCellStatusSuccess Indicates a successful or positive action
const TableCellStatusSuccess = "success"

//TableCellStatusInfo Indicates a neutral informative change or action
const TableCellStatusInfo = "info"

//TableCellStatusWarning Indicates a warning that might need attention
const TableCellStatusWarning = "warning"

//TableCellStatusDanger Indicates a dangerous or potentially negative action
const TableCellStatusDanger = "danger"

//TableRowStatusActive Applies the hover color to the table row or table cell
const TableRowStatusActive = TableCellStatusActive

//TableRowStatusSuccess Indicates a successful or positive action
const TableRowStatusSuccess = TableCellStatusSuccess

//TableRowStatusInfo Indicates a neutral informative change or action
const TableRowStatusInfo = TableCellStatusInfo

//TableRowStatusWarning Indicates a warning that might need attention
const TableRowStatusWarning = TableCellStatusWarning

//TableRowStatusDanger Indicates a dangerous or potentially negative action
const TableRowStatusDanger = TableCellStatusDanger

//TableHoverRows adds a hover effect (grey background color) on table rows
const TableHoverRows = "table-hover"

//TableBorderedTable adds borders on all sides of the table and cells
const TableBorderedTable = "table-bordered"

//TableStripedRows adds zebra-stripes to a tabl
const TableStripedRows = "table-striped"

//TableCondensedTable makes a table more compact by cutting cell padding in half
const TableCondensedTable = "table-condensed"

//TableResponsiveTable creates a responsive table. The table will then scroll horizontally on small devices (under 768px). When viewing on anything larger than 768px wide, there is no difference
const TableResponsiveTable = "table-responsive"
