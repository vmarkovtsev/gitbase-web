import React, { Component, Fragment } from 'react';
import PropTypes from 'prop-types';
import ReactTable from 'react-table';
import { Button } from 'react-bootstrap';
import 'react-table/react-table.css';
import './ResultsTable.less';

class ResultsTable extends Component {
  render() {
    const { showCode, showUAST } = this.props;
    const columns = this.props.response.meta.headers.map(col => ({
      Header: col,
      id: col,
      accessor: row => {
        let v = row[col];

        // Array of hashes to string
        if (Array.isArray(v) && v.every(e => typeof e === 'string')) {
          v = v.join('\n');
        }

        switch (typeof v) {
          case 'boolean':
            return v.toString();
          case 'string':
            // assume any multiline string is code
            if (v.indexOf('\n') > -1) {
              return (
                <Button
                  bsStyle="gbpl-tertiary"
                  className="btn-compact"
                  onClick={() => showCode(v)}
                >
                  CODE
                </Button>
              );
            }
            return v;
          case 'object':
            // UAST column
            const protobufs = row[`__${col}-protobufs`];
            return (
              <Button
                bsStyle="gbpl-tertiary"
                className="btn-compact"
                onClick={() => showUAST(v, protobufs)}
              >
                UAST
              </Button>
            );
          default:
            return v;
        }
      }
    }));

    return (
      <Fragment>
        <ReactTable
          className="results-table"
          data={this.props.response.data}
          columns={columns}
          defaultPageSize={10}
          minRows={0}
        />
        {this.props.response.meta.limit && (
          <div className="limit-text">
            (results have a forced LIMIT of {this.props.response.meta.limit})
          </div>
        )}
      </Fragment>
    );
  }
}

ResultsTable.propTypes = {
  // Must be a success response
  response: PropTypes.object.isRequired,
  showCode: PropTypes.func.isRequired,
  showUAST: PropTypes.func.isRequired
};

export default ResultsTable;
