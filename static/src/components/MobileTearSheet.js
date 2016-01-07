import React from 'react';
import {StylePropable} from 'material-ui/lib/mixins';

const MobileTearSheet = React.createClass({

  propTypes: {
    children: React.PropTypes.node,
    height: React.PropTypes.number,
  },

  contextTypes: {
    muiTheme: React.PropTypes.object,
  },

  getDefaultProps() {
    return {
      height: 500,
    };
  },

  render() {
    const styles = {
      root: {
        float: 'left',
        marginBottom: 24,
        marginRight: 24,
        width: 200,
      },
      container: {
        border: 'solid 1px #d9d9d9',
        borderBottom: 'none',
        height: this.props.height,
        overflow: 'hidden',
      },
      bottomTear: {
        display: 'block',
        position: 'relative',
        marginTop: -2,
        width: 200,
      },
    };

    return (
      <div style={styles.root}>
        <div style={styles.container}>
          {this.props.children}
        </div>
        <img style={styles.bottomTear} src="bottom-tear.svg" />
      </div>
    );
  },

});

export default MobileTearSheet;
