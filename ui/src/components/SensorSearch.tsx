import { createStyles, Divider, LinearProgress, makeStyles, Theme } from '@material-ui/core';
import React, { useState } from 'react'
import SearchBar from './SearchBar'
import * as API from "../api";
import DataTable, { DataTableRow } from './DataTable';

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

export default function SensorSearch() {

    const classes = useStyles();

    /**------------ */

    const [pagination, setPagination] = useState(null as API.Pagination)
    const [sensorsRows, setSensorsRows] = useState(null as DataTableRow[])
    const [searchQuery, setSearchQuery] = useState(null)
    const [searchLoading, setSearchLoading] = useState(false)
    const doSearch = (query: string, page: number) => {
        setSearchLoading(true)
        API.searchSensors(query, page).then(
            res => {

                setPagination(res.pagination)
                let tableRows: DataTableRow[] = new Array()
                for (let row of res.rows) {
                    tableRows.push({
                        id: row.channel_id,
                        title: row.name,
                        subtitle: row.channel_name,
                    });
                }
                setSensorsRows(tableRows)
            },
            err => { console.error(err); }
        ).finally(() => {
            setSearchLoading(false)
        })
    }

    const handleSearchSubmit = (e: any) => {
        e.preventDefault();
        doSearch(searchQuery, 1);
    }

    /**------------ */

    const handleChangePage = (e: React.MouseEvent<HTMLButtonElement>, page: number) => {
        e.preventDefault();
        doSearch(searchQuery, page + 1); // We need to increase it by one as the TablePagintation component starts with zero
    }

    /**------------ */

    return (
        <div className={classes.root} >
            <SearchBar
                onChange={(e) => { setSearchQuery(e.target.value); }}
                onSubmit={handleSearchSubmit}
                label="Search"
                placeholder="Search a sensor name" />
            { searchLoading && <LinearProgress />}
            <Divider />
            { sensorsRows &&
                <DataTable rows={sensorsRows} pagination={pagination} onChangePage={handleChangePage} />}
            { searchLoading && sensorsRows && <LinearProgress />}
        </div>
    )
}
