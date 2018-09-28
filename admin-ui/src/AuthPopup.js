import React, { Component } from 'react';
import { Table } from 'reactstrap';
import {Modal,ModalHeader,ModalFooter,ModalBody,Button} from 'reactstrap';
class AuthPopup extends Component {
    constructor(props) {
        super(props);
        
        this.state = {
          IsDropping: "",
          Nodes: [],
          BlockHeights: {},
          IsCreating: false,
          PendingKeys: []
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

    approve(key) {
        fetch("/api/nodes/auth/" + this.props.nodeName + "/" + key + "/1")
      .then(res => res.json())
      .then(res => {
        this.reloadDetails();
      });
    }

    decline(key) {
        fetch("/api/nodes/auth/" + this.props.nodeName + "/" + key + "/0")
      .then(res => res.json())
      .then(res => {
        this.reloadDetails();
      });
    }


    reloadDetails() {
    fetch("/api/nodes/pendingauth/" + this.props.nodeName)
      .then(res => res.json())
      .then(res => {
        this.setState({PendingKeys:res})
      });
    }

    render() {
        var pendingKeys = this.state.PendingKeys.map((k) => {
            return (
                <tr>
                    <td>{k.substring(0,6)}...{k.substring(60)}</td>
                    <td>
                        <Button onClick={((e) => this.approve(k))}>Approve</Button>{' '}
                        <Button onClick={((e) => this.decline(k))}>Decline</Button>
                    </td>
                </tr>
            )
        })

        return (<Modal isOpen={this.props.isOpen} className={this.props.className}>
            <ModalHeader>Pending authorizations for {this.props.nodeName}</ModalHeader>
            <ModalBody>
            <Table striped>
            <thead>
                <tr>
                <th>Name</th>
                <th>&nbsp;</th>
                </tr>
            </thead>
            <tbody>
                {pendingKeys}
            </tbody>
            </Table>
            </ModalBody>
            <ModalFooter>
              <Button color="primary" onClick={this.props.onClose}>Done</Button>
            </ModalFooter>
          </Modal>)
    }
}
export default AuthPopup;