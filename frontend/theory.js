new Vue({
    el: '#app',
    data: {
        loading: true,
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
        ]
    },
    created() {
        fetch('/api/theory')
            .then(response => response.json())
            .then(data => {
                this.results = data;
                this.loading = false;
            });
    },
    methods: {
        cityResults(cityName) {
            return this.results.filter(result => result.name === cityName);
        }
    }
});

// /api/theory
// /api/drive-manual
// /api/drive-auto


