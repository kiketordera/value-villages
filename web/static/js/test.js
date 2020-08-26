
//Register Navbar component
Vue.component('navbar-component', {
    template: '<nav class="navbar navbar-expand-md navbar-dark bg-dark mb-4">   <a class="navbar-brand" href="#">ACME INTL</a><button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation"> <span class="navbar-toggler-icon"></span> </button> </nav>',
  })
  
  //Register Page header component
  Vue.component('page-heading-component', {
    template: '<h1 class="text-center">{{heading}}</h1>',
    data: function() {
      return {
        heading: 'ACME Staff List'
      }
    }
  });
  //Register staff list component
  Vue.component('staff-list-component', {
    template: '<table class="table table-bordered"> <tbody>' +
      ' <tr v-for="staff in staffs"> <td>{{staff.name}}</td> <td>{{staff.email}}</td> <td>{{staff.role}}</td></tr>' +
      '</tbody></table>',
    data: function() {
      return {
        staffs: [{
          name: 'John Doe',
          email: 'John.doe@acme.org',
          role: 'Central Executive Officer'
        }, {
          name: 'Rebbecca Dan',
          email: 'rebbecca.dan@acme.org',
          role: 'Backend Developer'
        }, {
          name: 'Tope Joshua',
          email: 'tope.joshua@acme.org',
          role: 'Financial Analyst'
        }, {
          name: 'Alima Fatima',
          email: 'alima.fatima@acme.org',
          role: 'Deputy CTO'
        }, {
          name: 'Sikiru Oluwaseun',
          email: 'sikiru.oluwaseun@acme.org',
          role: 'Project Manager'
        }, {
          name: 'Larry Greg',
          email: 'larry.greg@acme.org',
          role: 'Senior Developer'
        }, {
          name: 'Inna Brown',
          email: 'inna.brown@acme.org',
          role: 'Community Manager'
        }, {
          name: 'Tunde Ogundipe',
          email: 'tunde.ogundipe@acme.org',
          role: 'Chief Technology Officer'
        }, {
          name: 'Bald Kuma',
          email: 'bald.kuma@acme.org',
          role: 'Human Resource'
        }, {
          name: 'Ramon Aduragbemi',
          email: 'ramon.aduragbemi@acme.org',
          role: 'System Administrator'
        }, ]
      }
    },
  
  });
  
  //Root Instance
  new Vue({
    el: '#app',
    data: {},
  });
  
  