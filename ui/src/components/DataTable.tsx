import React from 'react';
import { makeStyles, useTheme, Theme, createStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableFooter from '@material-ui/core/TableFooter';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import FirstPageIcon from '@material-ui/icons/FirstPage';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import LastPageIcon from '@material-ui/icons/LastPage';

import * as API from "../api";
import { Divider } from '@material-ui/core';

/*----------------------- */
const useStyles1 = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            flexShrink: 0,
            marginLeft: theme.spacing(2.5),
        },
    }),
);

interface TablePaginationActionsProps {
    count: number;
    page: number;
    rowsPerPage: number;
    onChangePage: (event: React.MouseEvent<HTMLButtonElement>, newPage: number) => void;
    className?: string;
}

function TablePaginationActions(props: TablePaginationActionsProps) {
    const classes = useStyles1();
    const theme = useTheme();
    const { count, page, rowsPerPage, onChangePage } = props;

    const handleFirstPageButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onChangePage(event, 0);
    };

    const handleBackButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onChangePage(event, page - 1);
    };

    const handleNextButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onChangePage(event, page + 1);
    };

    const handleLastPageButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onChangePage(event, Math.max(0, Math.ceil(count / rowsPerPage) - 1));
    };

    return (
        <div className={classes.root}>
            <IconButton
                onClick={handleFirstPageButtonClick}
                disabled={page === 0}
                aria-label="first page"
            >
                {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
            </IconButton>
            <IconButton onClick={handleBackButtonClick} disabled={page === 0} aria-label="previous page">
                {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
            </IconButton>
            <IconButton
                onClick={handleNextButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="next page"
            >
                {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
            </IconButton>
            <IconButton
                onClick={handleLastPageButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="last page"
            >
                {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
            </IconButton>
        </div>
    );
}

/*----------------------- */

const useStyles2 = makeStyles({
    table: {
        minWidth: 500,
    },
    tableRow: {
        cursor: 'pointer',
        '&:hover': {
            background: '#b4b4b433',
        }
    }
});

/**-------- */

export type DataTableRow = {
    id?: number;
    title?: string;
    subtitle?: string;
    onClick?: (event: React.MouseEvent<HTMLButtonElement>, id: number) => void;
}

interface Props {
    pagination: API.Pagination;
    rows: DataTableRow[];
    onChangePage: (event: React.MouseEvent<HTMLButtonElement>, newPage: number) => void;
}

/**-------- */
export default function DataTable(props: Props) {
    const classes = useStyles2();

    return (
        <TableContainer component={Paper} >
            <Table className={classes.table} aria-label="custom pagination table">
                <TableBody>
                    {props?.rows && props?.rows?.map((row, index) => (
                        <TableRow key={index} className={classes.tableRow}>
                            <TableCell component="th" scope="row">
                                {row?.title + "  - "}
                                <span style={{ color: "#888" }}>
                                    {row?.subtitle}
                                </span>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
                <TableFooter>
                    <TableRow >
                        <TablePagination
                            // rowsPerPageOptions={[5, 10, 25, { label: 'All', value: -1 }]}
                            rowsPerPageOptions={[]}
                            colSpan={3}
                            count={props?.pagination?.total_entries}
                            // count={-1}
                            rowsPerPage={props?.pagination?.total_entries > 200 ? 200 : props?.pagination?.total_entries}
                            page={props?.pagination?.current_page - 1}
                            // SelectProps={{
                            //     inputProps: { 'aria-label': 'rows per page' },
                            //     native: true,
                            // }}
                            onChangePage={props.onChangePage}
                            // onChangeRowsPerPage={() => { }}
                            ActionsComponent={TablePaginationActions}
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}
