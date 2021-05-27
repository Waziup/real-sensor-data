import { Box, createStyles, Divider, LinearProgress, makeStyles, Theme } from '@material-ui/core';
import React, { useEffect, useState } from 'react'
import SearchBar from '../components/SearchBar'
import * as API from "../api";
import DataTable, { DataTableRow } from '../components/DataTable';
import LoginForm from './LoginForm';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            padding: '5%',
            textAlign: 'center',
        },
        table: {
            marginTop: '10px',
        }
    }),
);

/**------------ */

interface Props {
    onSearchResultClick?: (dataRow: any) => void;
}

/**------------ */

export default function MyPushSensors(props: Props) {

    const classes = useStyles();

    /**------------ */

    useEffect(() => {
        load(1);
        // return () => {}
    }, [])

    /**------------ */

    const [err, setErr] = useState<API.HttpError>(null)
    const [pagination, setPagination] = useState(null as API.Pagination)
    const [sensorsRows, setSensorsRows] = useState(null as DataTableRow[])
    const [searchLoading, setLoading] = useState(false)
    const load = (page: number) => {
        setLoading(true)
        API.getMyPushSensors(page).then(
            res => {

                setPagination(res.pagination)
                let tableRows: DataTableRow[] = new Array()
                if (res.rows) {
                    for (let row of res.rows) {
                        tableRows.push({
                            id: row.id,
                            title: row.name,
                            subtitle: row.channel_name,
                        });
                    }
                }
                setSensorsRows(tableRows);
                setErr(null);
            },
            err => {
                console.error(err);
                setErr(err)
            }
        ).finally(() => {
            setLoading(false)
        })
    }

    /**------------ */

    const handleChangePage = (e: React.MouseEvent<HTMLButtonElement>, page: number) => {
        e.preventDefault();
        load(page + 1); // We need to increase it by one as the TablePagination component starts with zero
    }

    /**------------ */

    const handleTableRowClick = (e: any, rowIndex: number) => {
        if (props.onSearchResultClick) {
            props.onSearchResultClick(sensorsRows[rowIndex])
        }
    }

    /**------------ */

    let dataTable = null;
    if (sensorsRows) {
        if (sensorsRows.length) {
            dataTable = <DataTable rows={sensorsRows} pagination={pagination} onChangePage={handleChangePage}
                onRowClick={handleTableRowClick} />
        } else {
            dataTable = <Box component="span">No Sensors Found!</Box>
        }
    }

    /**------------ */

    // If authorization failed
    if (err && err.status == 401) {
        return <LoginForm onSuccess={() => { load(1); }} />
    }

    /**------------ */

    return (
        <div className={classes.root} >
            { searchLoading && <LinearProgress />}

            <Box component="div" mt={5}>
                {dataTable}
            </Box>

            { searchLoading && sensorsRows && <LinearProgress />}
        </div>
    )
}
