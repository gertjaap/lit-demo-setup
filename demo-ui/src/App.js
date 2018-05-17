import React, { Component } from 'react';

import { withStyles } from 'material-ui/styles';
import CssBaseline from 'material-ui/CssBaseline';
import './App.css';
import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import Grid from 'material-ui/Grid';
import IconButton from 'material-ui/IconButton';
import GavelIcon from '@material-ui/icons/Gavel';

import logo from './logo.svg';
const styles = theme => ({
  table: {
    minWidth: 700,
  },
});


class App extends Component {
  constructor(props) {
    super(props);
    
    this.state = {
      BlockHeight: 0,
    };
  }

  updateBlockHeight() {
    fetch("/api/chain/height")
    .then(res => res.json())
    .then(res => {
      this.setState({BlockHeight:res});
    })
  }


  mineBlock() {
    fetch("/api/chain/mine")
    .then(res => res.json())
    .then(res => {
      this.updateBlockHeight();
    })
  }


  render() {
    return (
       <div className="App">
        <CssBaseline />
        <AppBar position="static">
          <Toolbar>
            <Typography variant="title" color="inherit">
              LIT Demo Environment - Blockheight: {this.state.BlockHeight}
            </Typography>
            <IconButton onClick={() => this.mineBlock()} aria-label="Mine">
                      <GavelIcon />
            </IconButton>
          </Toolbar>
        </AppBar>
        <Grid container spacing={24}>
          <Grid xs={6}>
          <iframe title="Node 1" style={{width:'100%', height:'900px'}} src="/api/redirecttowebui?host=localhost&port=51001"></iframe>
          </Grid>
          <Grid xs={6}>
          <iframe title="Node 2" style={{width:'100%', height:'900px'}} src="/api/redirecttowebui?host=localhost&port=51002&alternative=1"></iframe>
          </Grid>
        </Grid>
      </div>
    );
  }

  componentDidMount() {
    this.updateBlockHeight();
    setInterval(() => {this.updateBlockHeight()}, 10000);
  }
}

export default withStyles(styles)(App);

