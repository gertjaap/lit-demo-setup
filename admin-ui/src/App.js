import React, { Component } from 'react';
import './App.css';
import { ScaleLoader } from 'react-spinners';
import { Navbar, NavbarBrand } from 'reactstrap';
import AuthPopup from './AuthPopup';
import LogsPopup from './LogsPopup';
import {Row,Col} from 'reactstrap';
import { Table } from 'reactstrap';
import { Button } from 'reactstrap';

class App extends Component {

  constructor(props) {
    super(props);
    
    this.state = {
      IsDropping: "",
      Nodes: [],
      BlockHeights: {},
      IsCreating: false,
      authPopupOpen: false,
      authPopupNodeName: '',
      logsPopupOpen: false,
      logsPopupNodeName: ''
    };

    this.newNode = this.newNode.bind(this);
    this.updateBlockHeight = this.updateBlockHeight.bind(this);
    this.update = this.update.bind(this);
    this.mineBlock = this.mineBlock.bind(this);
    this.dropNode = this.dropNode.bind(this);
    this.showNodeAuth = this.showNodeAuth.bind(this);
    this.restartNode = this.restartNode.bind(this);
    this.closeAuthPopup = this.closeAuthPopup.bind(this);
    this.closeLogsPopup = this.closeLogsPopup.bind(this);
  }

  updateBlockHeight() {
    fetch("/api/chain/height")
    .then(res => res.json())
    .then(res => {
      this.setState({BlockHeights:res});
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
    window.open('/api/redirecttowebui?host=' + document.location.hostname + '&port=' + node.PublicRpcPort)
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

  restartNode(node) {
    fetch("/api/nodes/restart/" + node.Name)
    .then(res => res.json())
    .then(res => {
       this.update();
    });
  }

  showNodeAuth(node) {
    this.setState({
      authPopupOpen: true,
      authPopupNodeName: node.Name
    })
  }

  closeAuthPopup() {
    this.setState({
      authPopupOpen: false,
      authPopupNodeName: ''
    })
  }

  showNodeLogs(node) {
    this.setState({
      logsPopupOpen: true,
      logsPopupNodeName: node.Name
    })
  }

  closeLogsPopup() {
    this.setState({
      logsPopupOpen: false,
      logsPopupNodeName: ''
    })
  }


  newNode() {
    this.setState({IsCreating:true})
    fetch("/api/nodes/new")
    .then(res => res.json())
    .then(res => {
    
      this.setState({ IsCreating: false})
      this.update();
    });
  }

  update() { 
    if(!this.state.IsCreating) {
      fetch("/api/nodes/list")
      .then(res => res.json())
      .then(res => {
        this.setState({Nodes:res})
      });
      this.updateBlockHeight();
    }
  }

  fundNode(node) {
    fetch("/api/nodes/fund/" + node.Name)
    .then(res => res.json())
    .then(res => {
      this.update();
    });
  }

  render() {
    let creation = null;
    let addNew = null;
    if(this.state.IsCreating) {
      creation = ( <tr>
          <td colspan="4"><ScaleLoader /> Creating new node...</td>
        </tr>)
    } else {
      addNew = ( <tr>
        <td colspan="4"><Button onClick={this.newNode}>Add new</Button></td>
      </tr>)
    }

    var blockHeights = Object.keys(this.state.BlockHeights).map((k) => {
      return <Col xs={4}><b>{k}:</b><br/><h1>{this.state.BlockHeights[k]}</h1></Col>;
    });

    var nodes = this.state.Nodes.map((n) => {
      return <tr>
        <td>{n.Name}</td>
        <td>{n.Address}</td>
        <td>{window.location.hostname}:{n.PublicLitPort}</td>
        <td>
          <Button onClick={((e) => { this.showNodeAuth(n); })}>Auth</Button>{' '}
          <Button onClick={((e) => { this.showNodeLogs(n); })}>Logs</Button>{' '}
          <Button onClick={((e) => { this.restartNode(n); })}>Restart</Button>{' '}
          <Button onClick={((e) => { this.dropNode(n); })}>Delete</Button>{' '}
        </td>
        </tr>;
    })

    return (
      <div className="App">
        
        <Navbar color="light" light expand="md">
          <NavbarBrand href="/">Lit Demo Environment</NavbarBrand>
        </Navbar>

        <Row>
          {blockHeights}
        </Row>

        <Table striped>
          <thead>
            <tr>
              <th>Name</th>
              <th>Status</th>
              <th>Funds</th>
              <th>&nbsp;</th>
            </tr>
          </thead>
          <tbody>
            {nodes}
            {creation}
            {addNew}
          </tbody>
        </Table>

        <AuthPopup isOpen={this.state.authPopupOpen} onClose={this.closeAuthPopup} nodeName={this.state.authPopupNodeName} />
        <LogsPopup isOpen={this.state.logsPopupOpen} onClose={this.closeLogsPopup} nodeName={this.state.logsPopupNodeName} />
      </div>
    );
  }

  componentDidMount() {
    this.update();
    setInterval(() => {this.update()}, 10000);
  }
}

export default App;

/**
 * 
 * 
 * <CssBaseline />
        <AppBar position="static">
          <Toolbar>
            <Typography variant="title" color="inherit">
              LIT Demo Environment - {blockHeights}
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
        
 */