new Vue({
    el: '#app',
    data: {
        loading: true,
        activeRoute: '',
        results: [],
        citiesOrder: [
            { name: "Kutaisi", order: 1 },
            { name: "Batumi", order: 2 },
            { name: "Telavi", order: 3 },
            { name: "Akhaltsikhe", order: 4 },
            { name: "Zugdidi", order: 5 },
            { name: "Gori", order: 6 },
            { name: "Poti", order: 7 },
            { name: "Ozurgeti", order: 8 },
            { name: "Sachkhere", order: 9 },
            { name: "Rustavi", order: 10 }
        ],
        currentRoute: '/api/get-auto',
        currentTitle: 'Drive city (auto)'
    },
    created() {
        this.fetchData();
    },
    methods: {
        cityResults(cityName) {
            return this.results.filter(result => result.name === cityName);
        },
        setCurrentRoute(route) {
            this.loading = true;
            this.activeRoute = route;
            this.currentRoute = route;

            if (route === '/api/get-auto') {
                this.currentTitle = 'Drive city (auto)';
            } else if (route === '/api/get-manual') {
                this.currentTitle = 'Manual Exam Dates';
            } else {
                this.currentTitle = 'Theory Exam Dates';
            }

            this.fetchData();
        },
        fetchData() {
            fetch(this.currentRoute)
                .then(response => response.json())
                .then(data => {
                    this.results = data;
                    this.loading = false;
                })
                .catch(error => {
                    console.error("There was an error fetching the data:", error);
                    this.loading = false;
                });
        }
    }
});
