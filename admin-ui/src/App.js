import React, { Component } from 'react';
import './App.css';
import { ScaleLoader } from 'react-spinners';
import { Navbar, NavbarBrand } from 'reactstrap';
import AuthPopup from './AuthPopup';
import QrPopup from './QrPopup';
import LogsPopup from './LogsPopup';
import {Row,Col} from 'reactstrap';
import { Table } from 'reactstrap';
import { Button } from 'reactstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHandshake, faQrcode } from '@fortawesome/free-solid-svg-icons'
import { faLock } from '@fortawesome/free-solid-svg-icons'
import { faFileAlt } from '@fortawesome/free-solid-svg-icons'
import { faSync } from '@fortawesome/free-solid-svg-icons'
import { faTrashAlt } from '@fortawesome/free-solid-svg-icons'

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
      logsPopupNodeName: '',
      qrPopupOpen: false,
      qrPopupNodeName: '',
      qrPopupNodeUrl: '',
      qrPopupMode: 'pair',
      qrPopupNodeAddress: ''
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
    this.closeQrPopup = this.closeQrPopup.bind(this);
    this.showNodeQr = this.showNodeQr.bind(this);
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

  closeQrPopup() {
    this.setState({
      qrPopupOpen: false,
      qrPopupNodeName: ''
    })
  }

  showNodeLogs(node) {
    this.setState({
      logsPopupOpen: true,
      logsPopupNodeName: node.Name
    })
  }

  showNodeQr(node) {
    this.setState({
      qrPopupOpen: true,
      qrPopupNodeName: node.Name,
      qrPopupNodeAddress: node.Address,
      qrPopupMode:'pair',
      qrPopupNodeUrl: node.Address + '@' + window.location.hostname + ':'  + node.PublicLitPort
    })
  }

  showNodeQrPay(node) {
    this.setState({
      qrPopupOpen: true,
      qrPopupNodeName: node.Name,
      qrPopupNodeAddress: node.Address,
      qrPopupMode:'pay',
      qrPopupNodeUrl: node.Address + '@' + window.location.hostname + ':'  + node.PublicLitPort
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
          <td colspan="8"><ScaleLoader /> Creating new node...</td>
        </tr>)
    } else {
      addNew = ( <tr>
        <td colspan="8"><Button onClick={this.newNode}>Add new</Button></td>
      </tr>)
    }

    var blockHeights = Object.keys(this.state.BlockHeights).map((k) => {
      return <Col xs={4}><b>{k}:</b><br/><h1>{this.state.BlockHeights[k]}</h1></Col>;
    });

    var nodes = this.state.Nodes.map((n) => {
      var balances = {};
      balances[257] = 0;
      balances[258] = 0;
      balances[262] = 0;
      if(n.Channels !== null) {
        n.Channels.forEach((c) => {
          balances[c.CoinType] += c.MyBalance / 100000000;
        });
      }

      return <tr>
        <td>{n.Name}</td>
        <td>{n.Address}  <Button title="Pay this node" onClick={((e) => { this.showNodeQrPay(n); })}><FontAwesomeIcon icon={faQrcode} /></Button></td>
        <td>{window.location.hostname}:{n.PublicLitPort}</td>
        <td><pre>{balances[257]}</pre></td>
        <td><pre>{balances[258]}</pre></td>
        <td><pre>{balances[262]}</pre></td>
        <td>{n.TrackerOK ? "OK" : "Fail"}</td>
        <td>
          <Button title="Pair with this node" onClick={((e) => { this.showNodeQr(n); })}><FontAwesomeIcon icon={faHandshake} /></Button>{' '}
          <Button title="Manage authorization requests" onClick={((e) => { this.showNodeAuth(n); })}><FontAwesomeIcon icon={faLock} /></Button>{' '}
          <Button title="Show logs of this node" onClick={((e) => { this.showNodeLogs(n); })}><FontAwesomeIcon icon={faFileAlt} /></Button>{' '}
          <Button onClick={((e) => { this.restartNode(n); })}><FontAwesomeIcon icon={faSync} /></Button>{' '}
          <Button onClick={((e) => { this.dropNode(n); })}><FontAwesomeIcon icon={faTrashAlt} /></Button>{' '}
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
              <th>Address</th>
              <th>Endpoint</th>
              <th>BTC</th>
              <th>LTC</th>
              <th>USD</th>
              <th>Tracker</th>
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
        <QrPopup isOpen={this.state.qrPopupOpen} onClose={this.closeQrPopup} mode={this.state.qrPopupMode} nodeAddress={this.state.qrPopupNodeAddress}  nodeUrl={this.state.qrPopupNodeUrl} nodeName={this.state.qrPopupNodeName} />
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