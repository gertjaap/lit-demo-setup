import React, { Component } from 'react';
import {Modal,ModalHeader,ModalFooter,ModalBody,Button} from 'reactstrap';
class LogsPopup extends Component {
    constructor(props) {
        super(props);
        
        this.state = {
          LogText: ""
        };

        this.reloadDetails = this.reloadDetails.bind(this);
    }

    componentDidMount() {
        if(this.props.nodeName !== '') {
            this.reloadDetails();
        }
    }

    componentDidUpdate(prevProps) {
        if(prevProps.nodeName !== this.props.nodeName || prevProps.isOpen === false && this.props.isOpen === true) {
            if(this.props.nodeName !== '') {
                this.reloadDetails();
            }
        }
    }

    reloadDetails() {
        fetch("/api/nodes/logs/" + this.props.nodeName)
        .then(res => res.text())
        .then(res => {
            this.setState({LogText:res})
        });
    }

    render() {

        return (<Modal isOpen={this.props.isOpen} className={this.props.className} style={{maxWidth:'100%', width:'90%'}}>
            <ModalHeader>Logs for {this.props.nodeName}</ModalHeader>
            <ModalBody>
                <div style={{whiteSpace: 'pre-wrap', fontFamily: 'monospace', width:'100%', height: '400px', overflowX:'auto', overflowY:'auto'}}>
                    {this.state.LogText}
                </div>
            </ModalBody>
            <ModalFooter>
            <Button color="secondary" onClick={this.props.onFullLog}>Full Log</Button>
            <Button color="primary" onClick={this.props.onClose}>Done</Button>
            </ModalFooter>
          </Modal>)
    }
}
export default LogsPopup;