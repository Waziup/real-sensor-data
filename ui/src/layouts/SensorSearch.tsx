import { Box, createStyles, Divider, LinearProgress, makeStyles, Theme } from '@material-ui/core';
import React, { useState } from 'react'
import SearchBar from '../components/SearchBar'
import * as API from "../api";
import DataTable, { DataTableRow } from '../components/DataTable';

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

export default function SensorSearch(props: Props) {

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
                if (res.rows) {
                    for (let row of res.rows) {
                        tableRows.push({
                            id: row.channel_id,
                            title: row.name,
                            subtitle: row.channel_name,
                        });
                    }
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

    return (
        <div className={classes.root} >
            <SearchBar
                onChange={(e) => { setSearchQuery(e.target.value); }}
                onSubmit={handleSearchSubmit}
                label="Search"
                placeholder="Search a sensor name" />
            { searchLoading && <LinearProgress />}

            <Box component="div" mt={5}>
                {dataTable}
            </Box>

            { searchLoading && sensorsRows && <LinearProgress />}
        </div>
    )
}
