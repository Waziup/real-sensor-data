import React, { useEffect, useState } from 'react';
import { makeStyles, useTheme, Theme, createStyles, withStyles } from '@material-ui/core/styles';
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
import Typography from 'material-ui/styles/typography';

import * as API from "../api";

/**---------------- */

import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { CircularProgress, Grid, TableHead } from '@material-ui/core';
let timeAgo: TimeAgo
try {
    TimeAgo.addDefaultLocale(en)
} catch (e) { //console.warn(e) 
}

/*----------------------- */

const StyledTableCell = withStyles((theme: Theme) =>
    createStyles({
        head: {
            backgroundColor: theme.palette.common.black,
            color: theme.palette.common.white,
        },
        body: {
            // fontSize: 18,
        },
    }),
)(TableCell);

const StyledTableRow = withStyles((theme: Theme) =>
    createStyles({
        root: {
            '&:nth-of-type(odd)': {
                backgroundColor: theme.palette.action.hover,
            },
        },
    }),
)(TableRow);


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
            background: '#edece6',
        }
    }
});

/**-------- */

interface Props {
    sensorId: number;
    userDevices: API.SensorType[];
    onRowClick?: (data: API.SensorPushSettings) => void;
}

/**-------- */
export default function DataTablePushSettings(props: Props) {
    const classes = useStyles2();

    var loadTimeout: NodeJS.Timeout = null
    useEffect(() => {
        timeAgo = new TimeAgo('en-US');
        loadPushSettings(1);
        return () => {
            clearTimeout(loadTimeout);
            loadTimeout = null;
        }
    }, [])

    /**-------------- */

    const handleRowClick = (event: React.MouseEvent<any>, index: number) => {
        if (props.onRowClick) props.onRowClick(pushSettingsList.rows[index])
    };

    /**-------------- */

    const intervalFormat = (value: number): string => {
        //value is in minutes

        let res = Math.floor(value / 1440);
        if (res) {
            return res + " day" + (res > 1 ? "s" : "");
        }

        res = Math.floor(value / 60);
        if (res) {
            return res + " hour" + (res > 1 ? "s" : "");
        }

        return value + " minute" + (value > 1 ? "s" : "");
    }

    /**--------------- */

    const handleChangePage = (e: React.MouseEvent<HTMLButtonElement>, page: number) => {
        e.preventDefault();
        loadPushSettings(page + 1); // We need to increase it by one as the TablePagination component starts with zero
    }
    /**--------------- */

    const [loadingPushSettings, setLoadingPushSettings] = useState<boolean>(false)
    const [pushSettingsList, setPushSettingsList] = useState<API.AllSensorPushSettings>(null)
    const loadPushSettings = (page: number) => {
        setLoadingPushSettings(true)
        API.getPushSettings(props.sensorId, page).then(
            res => {
                setPushSettingsList(res);
                loadTimeout = setTimeout(() => { loadPushSettings(page); }, 60 * 1000); // Refresh the table every minute
            },
            err => {
                console.error(err);
            }
        ).finally(() => setLoadingPushSettings(false))
    }

    /**--------------- */

    if (loadingPushSettings || props.userDevices === null || pushSettingsList === null) {
        return (<Grid container justify="center">
            <CircularProgress disableShrink />
        </Grid>)
    }

    /**--------------- */

    if (!pushSettingsList?.rows || pushSettingsList?.rows?.length == 0) {
        return (<Grid container justify="center">There is no push intervals configured for this sensor</Grid>)
    }

    /**--------------- */

    return (
        <TableContainer component={Paper} >
            <span style={{ fontWeight: "bold", fontSize: 18 }}>
                To edit each item, please click on it.
            </span>
            <Table className={classes.table} aria-label="custom pagination table">
                <TableHead>
                    <TableRow>
                        <StyledTableCell>Target Sensor</StyledTableCell>
                        <StyledTableCell>Status</StyledTableCell>
                        <StyledTableCell>Interval</StyledTableCell>
                        <StyledTableCell>Last Push</StyledTableCell>
                        <StyledTableCell>Total Push</StyledTableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {pushSettingsList?.rows && pushSettingsList?.rows?.map((row, index) => {

                        const devName = props.userDevices.find(o => o.devId == row.target_device_id)?.devName || row.target_device_id;
                        const sensorName = props.userDevices.find(o => o.id == row.target_sensor_id)?.name || row.target_sensor_id;

                        const lastPushTime = row.last_push_time ? timeAgo.format(new Date(row.last_push_time)) : "never";

                        return (<StyledTableRow key={index} className={classes.tableRow} onClick={(e) => handleRowClick(e, index)}>
                            <StyledTableCell component="th" scope="row" >
                                {devName + " - " + sensorName}
                            </StyledTableCell>
                            <StyledTableCell component="th" scope="row" >
                                {row.active ?
                                    <span style={{ color: "#44F" }}>Activate</span> :
                                    <span style={{ color: "#888" }}>Deactivate</span>
                                }
                            </StyledTableCell>
                            <StyledTableCell component="th" scope="row" >
                                {intervalFormat(row.push_interval)}
                            </StyledTableCell>
                            <StyledTableCell component="th" scope="row" >
                                {lastPushTime}
                            </StyledTableCell>
                            <StyledTableCell component="th" scope="row" >
                                {row?.pushed_count?.toLocaleString()}
                            </StyledTableCell>
                        </StyledTableRow>)
                    })}
                </TableBody>
                <TableFooter>
                    <TableRow >
                        <TablePagination
                            // rowsPerPageOptions={[5, 10, 25, { label: 'All', value: -1 }]}
                            rowsPerPageOptions={[]}
                            colSpan={10}
                            count={pushSettingsList?.pagination?.total_entries}
                            // count={-1}
                            rowsPerPage={pushSettingsList?.pagination?.total_entries > 200 ? 200 : pushSettingsList?.pagination?.total_entries}
                            page={pushSettingsList?.pagination?.current_page - 1}
                            // SelectProps={{
                            //     inputProps: { 'aria-label': 'rows per page' },
                            //     native: true,
                            // }}
                            onChangePage={handleChangePage}
                            // onChangeRowsPerPage={() => { }}
                            ActionsComponent={TablePaginationActions}
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}
