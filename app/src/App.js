import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
/*
[
    {
        "uuid": "41855398-6ef6-43bd-8abd-ce50ece25169",
        "firstName": "",
        "lastName": ""
    }
]
*/
class UsersTable extends Component {
  headers(keys) {
    return (
      <thead>
        <tr>
          { keys.map((key, idx) => (
            <th><abbr title={key}>{key}</abbr></th>
          ))}
          <th>del</th>
        </tr>
      </thead>
    )
  }

  body(keys, data) {
    return (
      <tbody>
        { data.map((record, idx) => (
          <tr>
            { keys.map((key, jdx) => (
               <th>{record[key]}</th>
            ))}
            <th><button className="button" onClick={this.props.onDelete}>Delete</button></th>
          </tr>
        ))}
      </tbody>  
    )
  }
  render() {
    const keys = Object.keys(this.props.data[0])
    return (
    <div>
      <table className="table">
      {this.headers(keys)} 
      {this.body(keys, this.props.data)} 
    </table>
    </div>
    )
  }
}

class UserForm extends Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.state = {
      firstName: "",
      lastName: ""
    }
  }
  handleChange(event) {
    this.setState({value: event.target.value});
  }
  render() {
    return (
      <div className="column">
        <input className="input" type="text" value={this.state.firstName} onChange={this.handleChange} placeholder="First Name"/>
        <input className="input" type="text" value={this.state.lastName} onChange={this.handleChange} placeholder="Last Name"/>
        <button className="button is-fullwidth" onClick={this.props.onClick}>Add User +</button>
      </div>
    )
  }
}

class App extends Component {
  constructor(props) {
    super(props);
    this.addUser = this.addUser.bind(this);
    this.deleteUser = this.deleteUser.bind(this);
    this.state = {
      usersList: [{}]
    }
  }

  getUsersList() {
    fetch("http://localhost:1323/users").then((res) => {
      return res.json()
    }).then((j) => {
      this.setState({
        usersList: j
      })
    })
  }

  addUser() {
    fetch("http://localhost:1323/users", {
      method: "POST"
    }).then((res) => {
      return res.json()
    }).then((j) => {
      this.getUsersList()
    })
  }

  deleteUser(e) {
    const id = e.target.value;
    console.log(id);
  }

  componentDidMount() {
    this.getUsersList()
  }

  render() {
    return (
      <div>
        <div className="columns">
          <div className="column">
            <UsersTable data={this.state.usersList} onDelete={this.deleteUser}/> 
          </div>
            <UserForm onClick={this.addUser} /> 
        </div>
      </div>
    );
  }
}

export default App;
