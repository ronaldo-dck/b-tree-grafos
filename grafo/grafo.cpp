#include <iostream>
#include <list>
#include <vector>
#include <iomanip>
#include <fstream>
#include <algorithm>
#include <cmath>
using namespace std;

class Vertice
{
public:
    float x;
    float y;
    string label;
    Vertice()
    {
        x = 0;
        y = 0;
        label = "";
    }
    Vertice(float x, float y, string label)
    {
        this->x = x;
        this->y = y;
        this->label = label;
    }
};

class Aresta
{
public:
    float custo;
    string label;
    string v1, v2;
    Aresta(float custo, string label, string v1, string v2)
    {
        this->custo = custo;
        this->label = label;
        this->v1 = v1;
        this->v2 = v2;
    }
};

class Grafo
{
private:
    vector<Vertice> vertices;
    vector<Aresta> arestas;
    vector<vector<float>> custos;

    void geraMatrizCustos()
    {
        vector<vector<float>> aux(size(vertices), vector<float>(size(vertices)));

        for (int i = 0; i < size(arestas); i++)
        {
            int linha, coluna;

            for (int j = 0; j < size(vertices); j++)
            {
                if (arestas[i].v1 == vertices[j].label)
                {
                    linha = j;
                }
                else if (arestas[i].v2 == vertices[j].label)
                {
                    coluna = j;
                }
            }

            aux[linha][coluna] = arestas[i].custo;
            aux[coluna][linha] = arestas[i].custo;
        }

        custos = aux;
    }

public:
    Grafo(/* args */){};
    Grafo(string nameFile)
    {
        ifstream file(nameFile, ios::in);
        int n;
        file >> n;
        string label;
        float x, y;
        for (int i = 0; i < n; i++)
        {
            file >> label >> x >> y;
            InserirVertice(x, y, label);
        }

        file >> n;
        string v1, v2;
        for (int i = 0; i < n; i++)
        {
            file >> label >> v1 >> v2;
            InserirAresta(label, v1, v2);
        }

        file.close();
    };

    void SalvarGrafo(string nameFile)
    {
        ofstream file(nameFile, ios::trunc | ios::out);

        file << size(vertices) << endl;
        for (int i = 0; i < size(vertices); i++)
        {
            file << vertices[i].label << " " << vertices[i].x << " " << vertices[i].y << endl;
        }
        file << size(arestas) << endl;
        for (int i = 0; i < size(arestas); i++)
        {
            file << arestas[i].label << " " << arestas[i].v1 << " " << arestas[i].v2 << endl;
        }

        file.close();
    }

    void InserirVertice(float x, float y, string label)
    {
        Vertice novoVertice(x, y, label);
        vertices.push_back(novoVertice);
        geraMatrizCustos();
    };

    void InserirAresta(string label, string v1, string v2)
    {

        Vertice A, B;

        for (int i = 0; i < size(vertices); i++)
        {
            if (v1 == vertices[i].label)
            {
                A = vertices[i];
            }
            else if (v2 == vertices[i].label)
            {
                B = vertices[i];
            }
        }
        if (A.label == "" || B.label == "")
        {
            cout << "Vértices não encontrados";
        }

        float custo = sqrt(pow(A.x - B.x, 2) + pow(A.y - B.y, 2));

        Aresta novaAresta(custo, label, v1, v2);
        arestas.push_back(novaAresta);
        geraMatrizCustos();
    }

    void RemoverAresta(string label)
    {
        int index;
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].label == label)
            {
                index = i;
                break;
            }
        }

        arestas.erase(arestas.begin() + index);
        geraMatrizCustos();
    }

    void RemoverVertice(string label)
    {
        int index;
        for (int i = 0; i < size(vertices); i++)
        {
            if (vertices[i].label == label)
            {
                index = i;
                break;
            }
        }

        vertices.erase(vertices.begin() + index);
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].v1 == label || arestas[i].v2 == label)
            {
                arestas.erase(arestas.begin() + i);
            }
        }
        geraMatrizCustos();
    }

    void PrintInfos()
    {
        cout << "lista de arestas" << endl;
        for (int i = 0; i < size(arestas); i++)
        {
            cout << arestas[i].label << ": " << arestas[i].v1 << "->" << arestas[i].v2 << " custo: " << arestas[i].custo << endl;
        }

        cout << "lista de vertices" << endl;
        for (int i = 0; i < size(vertices); i++)
        {
            cout << vertices[i].label << ": (" << vertices[i].x << ", " << vertices[i].y << ")\n";
        }

        cout << "Matriz de custos" << endl
             << "     ";
        for (int i = 0; i < size(custos); i++)
        {
            cout << " " << vertices[i].label;
        }
        cout << endl;
        for (int i = 0; i < size(custos); i++)
        {
            cout << vertices[i].label << " :";
            for (int j = 0; j < size(custos); j++)
            {
                cout << fixed << setprecision(2) << custos[i][j] << " ";
            }
            cout << endl;
        }
    }

    bool Euleriano() {
        
    }
};

int main()
{
    Grafo formiga("grafoIn.txt");
    int op;
    string label, A, B;

    while (1)
    {
        cout << " 0 - Sair\n 1 - Inserir Vértice\n 2 - Inserir Aresta\n 3 - Remover Vértice\n 4 - Remover Aresta\n 5 - Salvar em arquivo\n";
        cin >> op;
        switch (op)
        {
        case 0:
            return 0;
            break;
        case 1:
            float x, y;
            cout << "Label: ";
            cin >> label;
            cout << "X: ";
            cin >> x;
            cout << "Y: ";
            cin >> y;

            formiga.InserirVertice(x, y, label);
            break;
        case 2:
            cout << "Label: ";
            cin >> label;
            cout << "Vertice A: ";
            cin >> A;
            cout << "Vertice B: ";
            cin >> B;

            formiga.InserirAresta(label, A, B);
            break;
        case 3:
            cout << "Label: ";
            cin >> label;
            formiga.RemoverVertice(label);
            break;
        case 4:
            cout << "Label: ";
            cin >> label;
            formiga.RemoverAresta(label);
            break;
        case 5:
            // cout << "Nome do arquivo: ";
            // cin >> label;
            formiga.SalvarGrafo("grafoOut.txt");
            break;
        case 9:
            formiga.PrintInfos();
            break;
        default:
            break;
        }
    }

    return 0;
}