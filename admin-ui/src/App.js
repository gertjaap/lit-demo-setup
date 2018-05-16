import React, { Component } from 'react';
import CssBaseline from 'material-ui/CssBaseline';
import { withStyles } from 'material-ui/styles';
import './App.css';
import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import Grid from 'material-ui/Grid';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import AttachMoneyIcon from '@material-ui/icons/AttachMoney';
import OpenInNewIcon from '@material-ui/icons/OpenInNew';
import ContentCopy from 'material-ui-icons/ContentCopy';
import GavelIcon from '@material-ui/icons/Gavel';
import copy from 'copy-to-clipboard';
import { ScaleLoader } from 'react-spinners';
const styles = theme => ({
  table: {
    minWidth: 700,
  },
});


class App extends Component {

  constructor(props) {
    super(props);
    
    this.state = {
      IsDropping: "",
      Nodes: [],
      BlockHeight: 0,
      IsCreating: false
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

  openLitUIForNode(node) {
    window.open('/api/redirecttolitwebui?host=' + document.location.hostname + '&port=' + node.PublicRpcPort)
  }

  dropNode(node) {
    this.setState({IsDropping:node.Name});
    fetch("/api/nodes/delete/" + node.Name)
    .then(res => res.json())
    .then(res => {
      this.setState({IsDropping: ""});
      this.update();
    });
  }

  newNode() {
    this.setState({IsCreating:true})
    fetch("/api/nodes/new")
    .then(res => res.json())
    .then(res => {
      var nodes = this.state.Nodes;
      nodes.push(res);
      this.setState({Nodes:nodes, IsCreating: false})
    });
  }

  update() { 
    fetch("/api/nodes/list")
    .then(res => res.json())
    .then(res => {
      this.setState({Nodes:res})
    });
    this.updateBlockHeight();
  }

  fundNode(node) {
    fetch("/api/nodes/fund/" + node.Name)
    .then(res => res.json())
    .then(res => {
      this.update();
    });
  }

  copyPubNodeToClipboard(n) {
    copy(n.Address + '@' + document.location.hostname + ':' + n.PublicLitPort)
  }

  copyPrivNodeToClipboard(n) {
    copy(n.Address + '@' + n.Name)
  }


  render() {
    let creation = null;
    if(this.state.IsCreating) {
      creation = ( <Grid key="IsLoading" item xs={12}>
        <ScaleLoader />
        <Typography>
          Creating your new LIT node...
        </Typography>
      </Grid> )
    }
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
          {this.state.Nodes.map(n => {
            return (
              <Grid key={n.Name} item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="headline" component="h1" color="textSecondary">
                      {n.Name}
                    </Typography>
                    <Typography component="small">
                    <b>Public node address:</b><IconButton onClick={() => this.copyPubNodeToClipboard(n)} aria-label="Copy">
                      <ContentCopy />
                    </IconButton><br/>
                      {n.Address + '@' + document.location.hostname + ':' + n.PublicLitPort}<br/>
                    </Typography>
                    <Typography component="small">
                    <b>Private node address:</b><IconButton onClick={() => this.copyPrivNodeToClipboard(n)} aria-label="Copy">
                      <ContentCopy />
                    </IconButton><br/>
                    {n.Address + '@' + n.Name}
                    </Typography>
                    <Typography component="small">
                    <b>Balance:</b><br/>
                    {(n.Balances[257] / 100000) + ' mBTC-R'}
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <IconButton aria-label="Delete" onClick={() => this.dropNode(n)}>
                      <DeleteIcon  />
                    </IconButton>
                    <IconButton aria-label="Fund" onClick={() => this.fundNode(n)}>
                      <AttachMoneyIcon />
                    </IconButton>
                    <IconButton aria-label="Open LIT UI"  onClick={() => this.openLitUIForNode(n)}>
                      <OpenInNewIcon />
                    </IconButton>
                  </CardActions>
                </Card>
              </Grid>
            );
          })}
          
          <Grid key="New" item xs={12} sm={6}>
            {creation}
            <Button onClick={this.newNode.bind(this)}  disabled={this.state.IsCreating} variant="fab" color="primary" aria-label="add">
              <AddIcon />
            </Button>
          </Grid>
        </Grid>

      </div>
    );
  }

  componentDidMount() {
    this.update();
    setInterval(() => {this.update()}, 10000);
  }
}

export default withStyles(styles)(App);
